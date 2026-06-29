package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	mathrand "math/rand"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
	scrapligoffi "github.com/scrapli/scrapligo/v2/ffi"
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

func TestConcurrency(t *testing.T) {
	dumboCmd := exec.CommandContext(
		t.Context(),
		"go",
		"run",
		"main.go",
	)

	dumboCmd.Dir = "../build/dummy_ssh_server"

	err := dumboCmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	// since we just do "start" give the server time to run (esp since we are doing go run so has
	// to build the thing too)
	time.Sleep(250 * time.Millisecond)

	t.Cleanup(
		func() {
			t.Log("cleanup dummy server")

			_ = dumboCmd.Process.Kill()
			_ = dumboCmd.Wait()
		},
	)

	for _, transportName := range []string{
		"bin",
		"ssh2",
	} {
		testName := fmt.Sprintf("concurrencty-%s", transportName)

		t.Run(testName, func(t *testing.T) {
			t.Logf("%s: starting", testName)
			defer t.Logf("%s: complete", testName)

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

			for range 200 {
				wg.Go(
					func() {
						// tiny sleep seems to make the test way more consistent -- at least locally
						// on darwin i think we get starved for ptys and weird shit happens w/out
						// this.
						time.Sleep(time.Duration(mathrand.Intn(100)) * time.Millisecond)

						c, err := scrapligocli.NewCli( //nolint: contextcheck
							"localhost",
							opts...,
						)
						if err != nil {
							t.Fatal(err)
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
