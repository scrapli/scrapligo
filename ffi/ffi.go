package ffi

import (
	"bytes"
	"context"
	"debug/elf"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/ebitengine/purego"
	scrapligoconstants "github.com/scrapli/scrapligo/v2/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
	scrapligologging "github.com/scrapli/scrapligo/v2/logging"
	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
)

var (
	mappingInst     *Mapping  //nolint: gochecknoglobals
	mappingInstOnce sync.Once //nolint: gochecknoglobals
)

const (
	darwin         = "darwin"
	linux          = "linux"
	libscrapliRepo = "https://github.com/scrapli/libscrapli"
)

// AssertNoLeaks is a dev/test type function that asserts (using the general purpose allocator used
// in the underlying libscrapli ffi layer) that there are no memory leaks.
func AssertNoLeaks() error {
	m, err := GetMapping()
	if err != nil {
		return scrapligoerrors.NewFfiError("failed asserting no memory leaks", err)
	}

	noLeaks := m.AssertNoLeaks()
	if noLeaks {
		return nil
	}

	return scrapligoerrors.NewFfiError("found memory leaks", nil)
}

func getZigStyleArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "aarch64"
	default:
		panic("unsupported arch")
	}
}

func isMuslFallback() bool {
	matches, _ := filepath.Glob("/lib/ld-musl-*.so.1")

	return len(matches) > 0
}

func isMusl() bool {
	f, err := elf.Open("/proc/self/exe")
	if err != nil {
		return isMuslFallback()
	}

	defer func() {
		_ = f.Close()
	}()

	for _, p := range f.Progs {
		if p.Type == elf.PT_INTERP {
			interp := make([]byte, p.Filesz)

			_, err = io.ReadFull(io.NewSectionReader(p, 0, int64(p.Filesz)), interp) //nolint: gosec
			if err != nil {
				return isMuslFallback()
			}

			return bytes.Contains(interp, []byte("ld-musl"))
		}
	}

	return isMuslFallback()
}

func getLibscrapliCachePath() string {
	overridePath := os.Getenv(scrapligoconstants.LibScrapliCacheOverrideEnv)
	if overridePath != "" {
		scrapligologging.Logger(
			scrapligologging.Info,
			"using libscrapli cache path override %q...",
			overridePath,
		)

		return overridePath
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = ".cache"
	}

	cacheDir = filepath.Join(cacheDir, "scrapli")

	scrapligologging.Logger(
		scrapligologging.Debug,
		"using libscrapli cache dir %q...",
		cacheDir,
	)

	return cacheDir
}

func getLibscrapliVersion() string {
	return scrapligoutil.GetEnvStrOrDefault(
		scrapligoconstants.LibScrapliVersionOverrideEnv,
		scrapligoconstants.LibScrapliVersion,
	)
}

func libscrapliVersionIsHash(version string) bool {
	return regexp.MustCompile(
		`(?i)^[0-9a-f]{7,40}$`,
	).MatchString(version)
}

func getLibscrapliAbi() string {
	if isMusl() {
		return "musl"
	}

	return "gnu"
}

func getLibscrapliTargetFilename(version string) string {
	var libFilename string

	switch runtime.GOOS {
	case darwin:
		libFilename = fmt.Sprintf(
			"libscrapli-%s-macos.%s.dylib",
			getZigStyleArch(),
			version,
		)
	case linux:
		libFilename = fmt.Sprintf(
			"libscrapli-%s-linux-%s.so.%s",
			getZigStyleArch(),
			getLibscrapliAbi(),
			version,
		)
	default:
		panic("unsupported platform")
	}

	return libFilename
}

func getLibscrapliTargetZigTriple() string {
	var zigTriple string

	switch runtime.GOOS {
	case darwin:
		zigTriple = fmt.Sprintf("%s-macos", getZigStyleArch())
	case linux:
		zigTriple = fmt.Sprintf("%s-linux-%s", getZigStyleArch(), getLibscrapliAbi())
	default:
		panic("unsupported platform")
	}

	return zigTriple
}

// EnsureLibscrapli ensures libscrapli is present at the cache path. It returns the final path
// or an error.
func EnsureLibscrapli(ctx context.Context) (string, error) {
	overridePath := os.Getenv(scrapligoconstants.LibScrapliPathOverrideEnv)
	if overridePath != "" {
		scrapligologging.Logger(
			scrapligologging.Info,
			"using libscrapli path override %q...",
			overridePath,
		)

		return overridePath, nil
	}

	version := getLibscrapliVersion()

	libFilename := getLibscrapliTargetFilename(version)

	cachePath := getLibscrapliCachePath()

	cachedLibFilename := fmt.Sprintf("%s/%s", cachePath, libFilename)

	scrapligologging.Logger(
		scrapligologging.Debug,
		"looking for libscrapli at %q...",
		cachedLibFilename,
	)

	_, err := os.Stat(cachedLibFilename)
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		return cachedLibFilename, err
	}

	scrapligologging.Logger(
		scrapligologging.Info,
		"libscrapli does not exist at %q, writing libscrapli to disk at that location...",
		cachedLibFilename,
	)

	err = writeLibScrapliToCache(ctx, version, cachedLibFilename)
	if err != nil {
		return "", err
	}

	return cachedLibFilename, nil
}

func writeHTTPContentsFromPath(
	ctx context.Context,
	path string,
	w io.Writer,
) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close() //nolint: errcheck

	if resp.StatusCode != http.StatusOK {
		return scrapligoerrors.NewFfiError(
			fmt.Sprintf(
				"non 200 status attempting to load content at '%s', status code: %d",
				path,
				resp.StatusCode,
			),
			nil,
		)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func writeLibScrapliToCache( //nolint: nonamedreturns
	ctx context.Context,
	version string,
	cachedLibFilename string,
) (err error) {
	err = os.MkdirAll(
		filepath.Dir(cachedLibFilename),
		scrapligoconstants.PermissionsOwnerReadWriteExecute,
	)
	if err != nil {
		panic("failed ensuring cache directory")
	}

	var f *os.File

	f, err = os.Create(cachedLibFilename) //nolint: gosec
	if err != nil {
		return err
	}

	defer func() {
		// best effort close and remove if errored
		_ = f.Close()

		if err == nil {
			return
		}

		_ = os.Remove(cachedLibFilename)
	}()

	zigTriple := getLibscrapliTargetZigTriple()

	releaseFilename := getLibscrapliTargetFilename(version)

	if libscrapliVersionIsHash(version) {
		scrapligologging.Logger(
			scrapligologging.Debug,
			"libscrapli target version is hash, attempting to build via docker...",
		)

		_, err = exec.LookPath("docker")
		if err != nil {
			return scrapligoerrors.NewFfiError("docker unavailable and version is a hash", err)
		}

		buildCmd := exec.CommandContext( //nolint: gosec
			ctx,
			"docker",
			"run",
			"-v",
			fmt.Sprintf("%s:/out", filepath.Dir(cachedLibFilename)),
			"-e",
			fmt.Sprintf("TAG=%s", version),
			"-e",
			fmt.Sprintf("TARGET=%s", zigTriple),
			"-e",
			fmt.Sprintf("OUT_NAME=%s", filepath.Base(f.Name())),
			"ghcr.io/scrapli/libscrapli/builder:latest",
		)

		err := buildCmd.Run()
		if err != nil {
			return scrapligoerrors.NewFfiError("failed running build container", err)
		}
	} else {
		scrapligologging.Logger(
			scrapligologging.Debug,
			"libscrapli target version is tag, attempting to fetch from github...",
		)

		err = writeHTTPContentsFromPath(
			ctx,
			fmt.Sprintf(
				"%s/releases/download/v%s/%s",
				libscrapliRepo,
				version,
				releaseFilename,
			),
			f,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetMapping returns the singleton Mapping instance that holds the bindings to the underlying
// libscrapli shared library.
func GetMapping() (*Mapping, error) {
	var onceErrorString string

	mappingInstOnce.Do(func() {
		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		libscrapliPath, err := EnsureLibscrapli(ctx)
		if err != nil {
			onceErrorString = err.Error()

			return
		}

		libScrapliFfi, err := purego.Dlopen(
			libscrapliPath,
			purego.RTLD_NOW|purego.RTLD_GLOBAL,
		)
		if err != nil {
			onceErrorString = fmt.Sprintf(
				"error loading libscrapli at file %q, err: %s",
				libscrapliPath,
				err,
			)

			return
		}

		mappingInst = &Mapping{
			Shared:  SharedMapping{},
			Session: SessionMapping{},
			Cli:     CliMapping{},
			Netconf: NetconfMapping{},
		}

		purego.RegisterLibFunc(&mappingInst.AssertNoLeaks, libScrapliFfi, "ls_assert_no_leaks")

		registerShared(mappingInst, libScrapliFfi)
		registerSession(mappingInst, libScrapliFfi)
		registerCli(mappingInst, libScrapliFfi)
		registerNetconf(mappingInst, libScrapliFfi)

		scrapligologging.Logger(
			scrapligologging.Debug,
			"took %fs to load ffi...",
			time.Since(start).Seconds(),
		)
	})

	if mappingInst == nil {
		if onceErrorString != "" {
			return nil, scrapligoerrors.NewFfiError(onceErrorString, nil)
		}

		return nil, scrapligoerrors.NewFfiError("ffi mapping unavailable", nil)
	}

	return mappingInst, nil
}
