package cli_test

import (
	"bytes"
	"context"
	"fmt"
	mathrand "math/rand"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
	scrapligoffi "github.com/scrapli/scrapligo/v2/ffi"
	scrapligologging "github.com/scrapli/scrapligo/v2/logging"
	scrapligooptions "github.com/scrapli/scrapligo/v2/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/v2/testhelper"
)

const (
	testHost = "localhost"
)

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

func TestConcurrency(t *testing.T) { //nolint: gocognit
	tmpDir := t.TempDir()

	dumboBin := fmt.Sprintf("%s/dumbo", tmpDir)

	dumboBuild := exec.CommandContext( //nolint: gosec
		t.Context(),
		"go",
		"build",
		"-o",
		dumboBin,
		"main.go",
	)

	dumboBuild.Dir = "../build/dummy_ssh_server"

	err := dumboBuild.Run()
	if err != nil {
		t.Fatal(err)
	}

	for _, transportName := range []string{
		"bin",
		"ssh2",
	} {
		testName := fmt.Sprintf("concurrency-%s", transportName)

		t.Run(testName, func(t *testing.T) {
			t.Logf("%s: starting", testName)
			defer t.Logf("%s: complete", testName)

			dumboCmd := exec.CommandContext(
				t.Context(),
				dumboBin,
			)

			err = dumboCmd.Start()
			if err != nil {
				t.Fatal(err)
			}

			time.Sleep(250 * time.Millisecond)

			t.Cleanup(
				func() {
					t.Log("cleanup dummy server")

					err = dumboCmd.Process.Kill()
					if err != nil {
						t.Fatal(err)
					}

					_ = dumboCmd.Wait()
				},
			)

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			wg := &sync.WaitGroup{}

			opts := []scrapligooptions.Option{
				scrapligooptions.WithPort(2222),
				scrapligooptions.WithUsername("admin"),
				scrapligooptions.WithPassword("password"),
			}

			if transportName == "bin" {
				opts = append(
					opts,
					scrapligooptions.WithTransportBin(),
					scrapligooptions.WithBinTransportExtraArgs("-F /dev/null"),
				)
			} else {
				opts = append(
					opts,
					scrapligooptions.WithTransportSSH2(),
				)
			}

			for range 500 {
				wg.Go(
					func() {
						// tiny sleep seems to make the test way more consistent -- at least locally
						// on darwin i think we get starved for ptys and weird shit happens w/out
						// this.
						time.Sleep(
							time.Duration(mathrand.Intn(100)) * time.Millisecond, //nolint:gosec
						)

						c, err := scrapligocli.NewCli( //nolint: contextcheck
							"localhost",
							opts...,
						)
						if err != nil {
							t.Error(err)

							return
						}

						_, err = c.Open(ctx)
						if err != nil {
							t.Fatal(err)
						}

						defer func() {
							_, _ = c.Close(ctx)
						}()

						r, err := c.SendInput(ctx, "show version")
						if err != nil {
							t.Fatal(err)
						}

						scrapligotesthelper.AssertEqual(t, false, r.Failed())
					},
				)
			}

			wg.Wait()
		})

		time.Sleep(time.Second)
	}
}

// TestLoggerCallbackLimit asserts that we can create more clis (with a logger configured) than
// purego's hard limit of 2,000 callbacks. prior to the logger dispatcher restructuring, every
// cli/netconf object with a logger configured allocated a purego callback at option apply time --
// those callback slots are never freed, so object number 2,001 would panic. with the dispatcher
// singleton there is exactly one purego callback for logging regardless of how many objects exist.
func TestLoggerCallbackLimit(t *testing.T) {
	// 2k max callbacks in purego
	const instanceCount = 2_001

	for i := range instanceCount {
		c, err := scrapligocli.NewCli(
			testHost,
			scrapligooptions.WithDefinitionFileOrName(scrapligocli.AristaEos),
			scrapligooptions.WithLogger(
				func(_ scrapligologging.LogLevel, _ string) {},
			),
		)
		if err != nil {
			t.Fatalf("failed creating cli instance %d: %v", i, err)
		}

		// GetOptions applies the options -- that is where logger (callback) registration happens,
		// and, historically, where the per-object purego callback was allocated; this lets us
		// exercise the callback path without needing to actually open a connection
		_, err = c.GetOptions()
		if err != nil {
			t.Fatalf("failed applying options for cli instance %d: %v", i, err)
		}
	}
}

func getCli(t *testing.T, f string) *scrapligocli.Cli {
	t.Helper()

	opts := []scrapligooptions.Option{
		scrapligooptions.WithUsername("admin"),
		scrapligooptions.WithPassword("admin"),
		scrapligooptions.WithLookupKeyValue("enable", "libscrapli"),
		scrapligooptions.WithDefinitionFileOrName(scrapligocli.AristaEos),
	}

	if *scrapligotesthelper.Record {
		opts = append(
			opts,
			scrapligooptions.WithPort(22022),
			scrapligooptions.WithSessionRecorderPath(f),
		)
	} else {
		opts = append(
			opts,
			scrapligooptions.WithTransportTest(),
			scrapligooptions.WithTestTransportF(f),
			scrapligooptions.WithReadSize(1),
		)
	}

	d, err := scrapligocli.NewCli(
		testHost,
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

	cleanedActual := scrapligotesthelper.CleanCliOutput(t, r.Result())

	testGoldenContent := scrapligotesthelper.ReadFile(t, testGoldenPath)

	if !bytes.Equal(cleanedActual, testGoldenContent) {
		scrapligotesthelper.FailOutput(t, cleanedActual, testGoldenContent)
	}

	scrapligotesthelper.AssertEqual(t, 22, r.Port)
	scrapligotesthelper.AssertEqual(t, testHost, r.Host)
	scrapligotesthelper.AssertNotDefault(t, r.StartTime)
	scrapligotesthelper.AssertNotDefault(t, r.EndTime())
	scrapligotesthelper.AssertNotDefault(t, r.ElapsedTimeSeconds)
	scrapligotesthelper.AssertNotDefault(t, r.Results)
	scrapligotesthelper.AssertNotDefault(t, r.ResultsRaw)
	scrapligotesthelper.AssertEqual(t, false, r.Failed())
}
