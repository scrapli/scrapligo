package cli_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

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
