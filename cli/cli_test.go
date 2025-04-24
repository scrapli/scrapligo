package cli_test

import (
	"bytes"
	"os"
	"testing"

	scrapligocli "github.com/scrapli/scrapligo/cli"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

const (
	testHost = "localhost"
)

func TestMain(m *testing.M) {
	scrapligotesthelper.Flags()

	os.Exit(m.Run())
}

func getDriver(t *testing.T, f string) *scrapligocli.Driver {
	t.Helper()

	opts := []scrapligooptions.Option{
		scrapligooptions.WithUsername("admin"),
		scrapligooptions.WithPassword("admin"),
		scrapligooptions.WithLookupKeyValue("enable", "libscrapli"),
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

	d, err := scrapligocli.NewDriver(
		scrapligocli.AristaEos,
		testHost,
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func closeDriver(t *testing.T, d *scrapligocli.Driver) {
	t.Helper()

	// we simply free since we dont record/care about any closing bits
	p, m := d.GetPtr()
	m.Shared.Free(p)
}

func assertResult(t *testing.T, r *scrapligocli.Result, testGoldenPath string) {
	t.Helper()

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
	scrapligotesthelper.AssertNotDefault(t, r.Host)
	scrapligotesthelper.AssertNotDefault(t, r.Results)
	scrapligotesthelper.AssertNotDefault(t, r.ResultsRaw)
	scrapligotesthelper.AssertEqual(t, false, r.Failed())
}
