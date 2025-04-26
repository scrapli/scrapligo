package netconf_test

import (
	"bytes"
	"os"
	"testing"

	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
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

func getNetconf(t *testing.T, f string) *scrapligonetconf.Netconf {
	t.Helper()

	opts := []scrapligooptions.Option{
		// note that netconf-admin bypasses enable secret stuff, without this was getting
		// permission denied committing things and such... but wanted to retain the enable
		// secret stuff since its nice to validate default mode gets acquired and stuff
		scrapligooptions.WithUsername("root"),
		scrapligooptions.WithPassword("password"),
		scrapligooptions.WithPort(830),
	}

	if *scrapligotesthelper.Record {
		opts = append(
			opts,
			scrapligooptions.WithSessionRecorderPath(f),
		)
	} else {
		opts = append(
			opts,
			scrapligooptions.WithTransportTest(),
			scrapligooptions.WithTestTransportF(f),
			scrapligooptions.WithReadSize(1),
			// see libscrapli notes in integration netconf tests
			scrapligooptions.WithOperationMaxSearchDepth(32),
		)
	}

	d, err := scrapligonetconf.NewNetconf(
		testHost,
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func closeNetconf(t *testing.T, d *scrapligonetconf.Netconf) {
	t.Helper()

	// we simply free since we dont record/care about any closing bits
	p, m := d.GetPtr()
	m.Shared.Free(p)
}

func assertResult(t *testing.T, r *scrapligonetconf.Result, testGoldenPath string) {
	t.Helper()

	cleanedActual := scrapligotesthelper.CleanNetconfOutput(t, r.Result)

	// we can't just write the cleaned stuff to disk because then chunk sizes will be wrong if we
	// just do the lazy cleanup method we are doing (and cant stop wont stop)
	testGoldenContent := scrapligotesthelper.ReadFile(t, testGoldenPath)
	cleanedGolden := scrapligotesthelper.CleanNetconfOutput(t, string(testGoldenContent))

	if !bytes.Equal(cleanedActual, cleanedGolden) {
		scrapligotesthelper.FailOutput(t, cleanedActual, cleanedGolden)
	}

	scrapligotesthelper.AssertEqual(t, r.Port, 830)
	scrapligotesthelper.AssertEqual(t, r.Host, testHost)
	scrapligotesthelper.AssertNotDefault(t, r.StartTime)
	scrapligotesthelper.AssertNotDefault(t, r.EndTime)
	scrapligotesthelper.AssertNotDefault(t, r.ElapsedTimeSeconds)
	scrapligotesthelper.AssertNotDefault(t, r.Host)
	scrapligotesthelper.AssertNotDefault(t, r.ResultRaw)
	scrapligotesthelper.AssertEqual(t, false, r.Failed)
}
