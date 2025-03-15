package netconf_test

import (
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
	opts := []scrapligooptions.Option{
		scrapligooptions.WithUsername("admin"),
		scrapligooptions.WithPassword("admin"),
		scrapligooptions.WithPort(22830),
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

func closeNetconf(t *testing.T, n *scrapligonetconf.Netconf, f string) {
	if *scrapligotesthelper.Record {
		p, m := n.GetPtr()
		m.Netconf.Free(p)

		return
	}

	n.Close()

	if !*scrapligotesthelper.Record {
		return
	}

	scrapligotesthelper.WriteFile(t, f, scrapligotesthelper.ReadFile(t, f))
}
