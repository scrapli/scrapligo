package ffi

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/ebitengine/purego"
	scrapligoassets "github.com/scrapli/scrapligo/assets"
	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligologging "github.com/scrapli/scrapligo/logging"
)

var (
	mappingInst     *Mapping  //nolint: gochecknoglobals
	mappingInstOnce sync.Once //nolint: gochecknoglobals
)

const (
	darwin = "darwin"
	linux  = "linux"
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

func getLibscrapliPath() (string, error) {
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

func writeLibScrapliToCache(cachedLibFilename string) error {
	var assetFilename string

	switch runtime.GOOS {
	case darwin:
		assetFilename = fmt.Sprintf(
			"lib/%s-macos/libscrapli.%s.dylib",
			getZigStyleArch(),
			scrapligoconstants.LibScrapliVersion,
		)
	case linux:
		abi := "gnu"

		if isMusl() {
			abi = "musl"
		}

		assetFilename = fmt.Sprintf(
			"lib/%s-linux-%s/libscrapli.so.%s",
			getZigStyleArch(),
			abi,
			scrapligoconstants.LibScrapliVersion,
		)
	default:
		panic("unsupported platform")
	}

	contents, err := scrapligoassets.Lib.ReadFile(assetFilename)
	if err != nil {
		return err
	}

	err = os.MkdirAll(
		filepath.Dir(cachedLibFilename),
		scrapligoconstants.PermissionsOwnerReadWriteExecute,
	)
	if err != nil {
		return err
	}

	err = os.WriteFile(cachedLibFilename, contents, scrapligoconstants.PermissionsOwnerReadWrite)
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

		libscrapliPath, err := getLibscrapliPath()
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

		// TODO can/should we only register ssh stuff if ssh and netconf if netconf? is the overhead
		// of the initial loading meaningful? this may be even more relevant/interesting if the
		// filesize and startup time matter (due to having to load from binary then write to disk)
		mappingInst = &Mapping{
			Driver:        DriverMapping{},
			DriverNetconf: DriverNetconfMapping{},
			Options: OptionMapping{
				Session:       SessionOptions{},
				Auth:          AuthOptions{},
				TransportBin:  TransportBinOptions{},
				TransportSSH2: TransportSSH2Options{},
			},
		}

		purego.RegisterLibFunc(&mappingInst.AssertNoLeaks, libScrapliFfi, "assertNoLeaks")

		registerDriver(mappingInst, libScrapliFfi)
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
