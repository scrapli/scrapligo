package cli_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"

	scrapligocli "github.com/scrapli/scrapligo/cli"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

var (
	isEosAvailable     *bool     //nolint: gochecknoglobals
	isEosAvailableOnce sync.Once //nolint: gochecknoglobals
)

func slowTests() []string {
	return []string{
		"send-input-enormous-srl-bin",
		"send-input-enormous-srl-ssh2",
	}
}

func eosAvailable() bool {
	isEosAvailableOnce.Do(func() {
		isAvailable := false

		o, err := exec.Command("docker", "ps").CombinedOutput()
		if err == nil {
			if bytes.Contains(o, []byte("ceos")) {
				isAvailable = true
			}
		}

		isEosAvailable = &isAvailable
	})

	return *isEosAvailable
}

func TestMain(m *testing.M) {
	scrapligotesthelper.Flags()

	exitCode := m.Run()

	if scrapligoffi.AssertNoLeaks() != nil {
		_, _ = fmt.Fprintln(os.Stderr, "memory leak(s) detected!")

		os.Exit(127)
	}

	_, _ = fmt.Fprintln(os.Stderr, "no memory leak(s) detected!")

	os.Exit(exitCode)
}

func getCli(t *testing.T, platform, transportName string) *scrapligocli.Cli {
	t.Helper()

	if strings.Contains(*scrapligotesthelper.Platforms, transportName) {
		t.Skipf("skipping platform %q, due to cli flag...", platform)
	}

	if strings.Contains(*scrapligotesthelper.Transports, transportName) {
		t.Skipf("skipping transport %q, due to cli flag...", transportName)
	}

	var host string

	opts := []scrapligooptions.Option{
		scrapligooptions.WithUsername("admin"),
	}

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
	case "telnet":
		opts = append(
			opts,
			scrapligooptions.WithTransportTelnet(),
		)
	default:
		t.Fatal("unsupported transport name")
	}

	if platform == scrapligocli.AristaEos.String() {
		if !eosAvailable() {
			t.Skip("skipping case, arista eos unavailable...")
		}

		opts = append(
			opts,
			scrapligooptions.WithPassword("admin"),
			scrapligooptions.WithLookupKeyValue("enable", "libscrapli"),
		)

		var port uint16

		if runtime.GOOS == "darwin" {
			host = "localhost"
			port = 22022
		} else {
			host = "172.20.20.17"
		}

		if transportName == "telnet" {
			port++
		}

		opts = append(
			opts,
			scrapligooptions.WithPort(port),
		)
	} else {
		opts = append(
			opts,
			scrapligooptions.WithPassword("NokiaSrl1!"),
		)

		if runtime.GOOS == "darwin" {
			host = "localhost"

			opts = append(
				opts,
				scrapligooptions.WithPort(21022),
			)
		} else {
			host = "172.20.20.16"
		}
	}

	d, err := scrapligocli.NewCli(
		platform,
		host,
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func assertResult(t *testing.T, r *scrapligocli.Result, testGoldenPath string) {
	t.Helper()

	if *scrapligotesthelper.Update {
		scrapligotesthelper.WriteFile(
			t,
			testGoldenPath,
			scrapligotesthelper.CleanCliOutput(t, r.Result()),
		)

		return
	}

	if !(*scrapligotesthelper.SkipSlow && slices.Contains(slowTests(), t.Name())) {
		cleanedActual := scrapligotesthelper.CleanCliOutput(t, r.Result())

		testGoldenContent := scrapligotesthelper.ReadFile(t, testGoldenPath)

		if !bytes.Equal(cleanedActual, testGoldenContent) {
			scrapligotesthelper.FailOutput(t, cleanedActual, testGoldenContent)
		}
	}

	scrapligotesthelper.AssertNotDefault(t, r.Port)
	scrapligotesthelper.AssertNotDefault(t, r.Host)
	scrapligotesthelper.AssertNotDefault(t, r.StartTime)
	scrapligotesthelper.AssertNotDefault(t, r.EndTime())
	scrapligotesthelper.AssertNotDefault(t, r.ElapsedTimeSeconds)
	scrapligotesthelper.AssertNotDefault(t, r.Host)
	scrapligotesthelper.AssertNotDefault(t, r.Results)
	scrapligotesthelper.AssertNotDefault(t, r.ResultsRaw)
	scrapligotesthelper.AssertEqual(t, false, r.Failed())
}
