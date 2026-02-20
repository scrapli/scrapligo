package netconf_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	scrapligocli "github.com/scrapli/scrapligo/cli"
	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

const localhost = "localhost"

func netconfTransports() []string {
	return []string{"bin", "ssh2"}
}

func TestMain(m *testing.M) {
	scrapligotesthelper.Flags()

	// ensure keys are chmod'd
	err := filepath.Walk(
		"./fixtures",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.Contains(filepath.Base(path), "key") {
				err = os.Chmod(path, 0o600)
				if err != nil {
					return fmt.Errorf("failed to chmod %s: %w", path, err)
				}
			}

			return nil
		},
	)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed ensuring key permissions")

		os.Exit(127)
	}

	exitCode := m.Run()

	if scrapligoffi.AssertNoLeaks() != nil {
		_, _ = fmt.Fprintln(os.Stderr, "memory leak(s) detected!")

		os.Exit(127)
	}

	_, _ = fmt.Fprintln(os.Stderr, "no memory leak(s) detected!")

	os.Exit(exitCode)
}

func getNetconf(t *testing.T, platform, transportName string) *scrapligonetconf.Netconf {
	t.Helper()

	if *scrapligotesthelper.Platforms != "all" &&
		!strings.Contains(*scrapligotesthelper.Platforms, platform) {
		t.Skipf("skipping platform %q, due to cli flag...", platform)
	}

	if *scrapligotesthelper.Platforms != "all" &&
		!strings.Contains(*scrapligotesthelper.Transports, transportName) {
		t.Skipf("skipping transport %q, due to cli flag...", transportName)
	}

	var opts []scrapligooptions.Option

	switch transportName {
	case "bin":
		opts = append(
			opts,
			scrapligooptions.WithTransportBin(),
		)
	case "ssh2":
		opts = append(
			opts,
			scrapligooptions.WithTransportSSH2(),
		)
	default:
		t.Fatal("unsupported transport name")
	}

	host := localhost

	switch platform {
	case scrapligocli.AristaEos.String():
		opts = append(
			opts,
			scrapligooptions.WithUsername("netconf-admin"),
			scrapligooptions.WithPassword("admin"),
		)

		if runtime.GOOS == scrapligoconstants.Darwin {
			opts = append(
				opts,
				scrapligooptions.WithPort(22830),
			)
		} else {
			host = "172.20.20.17"
		}
	case scrapligocli.NokiaSrlinux.String():
		opts = append(
			opts,
			scrapligooptions.WithUsername("admin"),
			scrapligooptions.WithPassword("NokiaSrl1!"),
		)

		if runtime.GOOS == scrapligoconstants.Darwin {
			opts = append(
				opts,
				scrapligooptions.WithPort(21830),
			)
		} else {
			host = "172.20.20.16"
		}
	default:
		// netopeer server
		opts = append(
			opts,
			scrapligooptions.WithUsername("root"),
			scrapligooptions.WithPassword("password"),
		)

		if runtime.GOOS == scrapligoconstants.Darwin {
			opts = append(
				opts,
				scrapligooptions.WithPort(23830),
			)
		} else {
			host = "172.20.20.18"
		}
	}

	n, err := scrapligonetconf.NewNetconf(
		host,
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return n
}

func xmlIsValid(result []byte) bool {
	decoder := xml.NewDecoder(bytes.NewReader(result))

	for {
		_, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				return true
			}

			return false
		}
	}
}

func assertResult(t *testing.T, r *scrapligonetconf.Result, testGoldenPath string) {
	t.Helper()

	if *scrapligotesthelper.Update {
		scrapligotesthelper.WriteFile(
			t,
			testGoldenPath,
			scrapligotesthelper.CleanNetconfOutput(t, r.Result),
		)

		return
	}

	cleanedActual := scrapligotesthelper.CleanNetconfOutput(t, r.Result)

	if !xmlIsValid(cleanedActual) {
		t.Fatal("result xml is invalid")
	}

	// skip golden comparo on a few tests due to output moving around
	if !slices.Contains(
		[]string{"TestGet/get-simple-netopeer-bin", "TestGet/get-simple-netopeer-ssh2"},
		t.Name(),
	) {
		// we can't just write the cleaned stuff to disk because then chunk sizes will be wrong if
		// we just do the lazy cleanup method we are doing (and cant stop wont stop)
		testGoldenContent := scrapligotesthelper.ReadFile(t, testGoldenPath)
		cleanedGolden := scrapligotesthelper.CleanNetconfOutput(t, string(testGoldenContent))

		if !bytes.Equal(cleanedActual, cleanedGolden) {
			scrapligotesthelper.FailOutput(t, cleanedActual, cleanedGolden)
		}
	}

	// we dont check failed since for now some things (cancel commit) fail expectedly, but we are
	// more just making sure the rpc was successful and we sent valid stuff etc.
	scrapligotesthelper.AssertNotDefault(t, r.StartTime)
	scrapligotesthelper.AssertNotDefault(t, r.EndTime)
	scrapligotesthelper.AssertNotDefault(t, r.ElapsedTimeSeconds)
	scrapligotesthelper.AssertNotDefault(t, r.Host)
	scrapligotesthelper.AssertNotDefault(t, r.ResultRaw)
}
