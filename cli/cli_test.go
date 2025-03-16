package cli_test

import (
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

func closeDriver(t *testing.T, d *scrapligocli.Driver, f string) {
	if *scrapligotesthelper.Record {
		p, m := d.GetPtr()
		m.Cli.Free(p)

		return
	}

	d.Close()

	if !*scrapligotesthelper.Record {
		return
	}

	scrapligotesthelper.WriteFile(t, f, scrapligotesthelper.ReadFile(t, f))
}
