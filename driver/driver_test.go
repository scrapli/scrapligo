package driver_test

import (
	scrapligodriver "github.com/scrapli/scrapligo/driver"
	"os"
	"testing"

	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestMain(m *testing.M) {
	scrapligotesthelper.Flags()

	os.Exit(m.Run())
}

func getDriver(t *testing.T, f string) *scrapligodriver.Driver {
	opts := []scrapligodriver.Option{
		scrapligodriver.WithUsername("admin"),
		scrapligodriver.WithPassword("admin"),
		scrapligodriver.WithLookupKeyValue("enable", "libscrapli"),
	}

	if *scrapligotesthelper.Record {
		opts = append(
			opts,
			scrapligodriver.WithPort(22022),
			scrapligodriver.WithSessionRecorderPath(f),
		)
	} else {
		opts = append(
			opts,
			scrapligodriver.WithTransportKind(scrapligodriver.TransportKindTest),
			scrapligodriver.WithTestTransportF(f),
			scrapligodriver.WithReadSize(1),
		)
	}

	d, err := scrapligodriver.NewDriver(
		string(scrapligodriver.AristaEos),
		"localhost",
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return d
}
