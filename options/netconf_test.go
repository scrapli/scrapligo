package options_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	scrapligonetconf "github.com/scrapli/scrapligo/v2/netconf"
	scrapligooptions "github.com/scrapli/scrapligo/v2/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/v2/testhelper"
)

func TestNETCONFOptions(t *testing.T) {
	d, err := scrapligonetconf.NewNetconf(
		"1.2.3.4",
		scrapligooptions.WithNetconfErrorTag("<errrrrror>"),
		scrapligooptions.WithNetconfPreferredVersion(scrapligooptions.NetconfVersion10),
		scrapligooptions.WithNetconfMessagePollIntervalNS(999999),
	)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := d.GetOptions()
	if err != nil {
		t.Fatal(err)
	}

	testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", "netconf.json"))
	if err != nil {
		t.Fatal(err)
	}

	if *scrapligotesthelper.Update {
		scrapligotesthelper.WriteFile(
			t,
			testGoldenPath,
			[]byte(actual),
		)

		return
	}

	testGoldenContent := string(scrapligotesthelper.ReadFile(t, testGoldenPath))

	if !strings.EqualFold(actual, testGoldenContent) {
		scrapligotesthelper.FailOutput(t, actual, testGoldenContent)
	}
}
