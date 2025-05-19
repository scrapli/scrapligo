package ffi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/ebitengine/purego"
	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligologging "github.com/scrapli/scrapligo/logging"
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

func isMusl() bool {
	_, err := os.Stat("/lib/ld-musl-x86_64.so.1")

	return err == nil
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

	var cacheDir string

	switch runtime.GOOS {
	case darwin:
		cacheDir = fmt.Sprintf("%s/Library/Caches/scrapli", os.Getenv(scrapligoconstants.HomeEnv))
	case linux:
		cacheDir = os.Getenv(scrapligoconstants.XdgCacheHomeEnv)
		if cacheDir == "" {
			cacheDir = fmt.Sprintf("%s/.cache/scrapli", os.Getenv(scrapligoconstants.HomeEnv))
		}
	default:
		panic("unsupported platform")
	}

	scrapligologging.Logger(
		scrapligologging.Debug,
		"using libscrapli cache dir %q...",
		cacheDir,
	)

	return cacheDir
}

// EnsureLibscrapli ensures libscrapli is present at the cache path. It returns the final path
// or an error.
func EnsureLibscrapli() (string, error) {
	overridePath := os.Getenv(scrapligoconstants.LibScrapliPathOverrideEnv)
	if overridePath != "" {
		scrapligologging.Logger(
			scrapligologging.Info,
			"using libscrapli path override %q...",
			overridePath,
		)

		return overridePath, nil
	}

	var libFilename string

	switch runtime.GOOS {
	case darwin:
		libFilename = fmt.Sprintf(
			"libscrapli.%s.dylib",
			scrapligoconstants.LibScrapliVersion,
		)
	case linux:
		libFilename = fmt.Sprintf(
			"libscrapli.so.%s",
			scrapligoconstants.LibScrapliVersion,
		)
	default:
		panic("unsupported platform")
	}

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

	err = writeLibScrapliToCache(cachedLibFilename)
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

	defer resp.Body.Close() //nolint

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

func writeLibScrapliToCache(cachedLibFilename string) error {
	f, err := os.Create(cachedLibFilename) //nolint: gosec
	if err != nil {
		return err
	}

	var releaseFilename string

	switch runtime.GOOS {
	case darwin:
		releaseFilename = fmt.Sprintf(
			"libscrapli-%s-macos.dylib.%s",
			getZigStyleArch(),
			scrapligoconstants.LibScrapliVersion,
		)
	case linux:
		abi := "gnu"

		if isMusl() {
			abi = "musl"
		}

		releaseFilename = fmt.Sprintf(
			"libscrapli-%s-linux-%s-%s.so",
			getZigStyleArch(),
			abi,
			scrapligoconstants.LibScrapliVersion,
		)
	default:
		panic("unsupported platform")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err = writeHTTPContentsFromPath(
		ctx,
		fmt.Sprintf(
			"%s/releases/download/v%s/%s",
			libscrapliRepo,
			scrapligoconstants.LibScrapliVersion,
			releaseFilename,
		),
		f,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetMapping returns the singleton Mapping instance that holds the bindings to the underlying
// libscrapli shared library.
func GetMapping() (*Mapping, error) {
	var onceErrorString string

	mappingInstOnce.Do(func() {
		start := time.Now()

		libscrapliPath, err := EnsureLibscrapli()
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
			Cli:     CliMapping{},
			Netconf: NetconfMapping{},
			Options: OptionMapping{
				Session:       SessionOptions{},
				Auth:          AuthOptions{},
				TransportBin:  TransportBinOptions{},
				TransportSSH2: TransportSSH2Options{},
			},
		}

		purego.RegisterLibFunc(&mappingInst.AssertNoLeaks, libScrapliFfi, "ls_assert_no_leaks")

		registerShared(mappingInst, libScrapliFfi)
		registerCli(mappingInst, libScrapliFfi)
		registerNetconf(mappingInst, libScrapliFfi)
		registerOptions(mappingInst, libScrapliFfi)

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
