package options_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	scrapligocli "github.com/kentik/scrapligo/v2/cli"
	scrapligooptions "github.com/kentik/scrapligo/v2/options"
	scrapligotesthelper "github.com/kentik/scrapligo/v2/testhelper"
)

func TestTransportBinOptions(t *testing.T) {
	d, err := scrapligocli.NewCli(
		"1.2.3.4",
		scrapligooptions.WithTransportBin(),
		scrapligooptions.WithBinTransportBinOverride("/myyyyy/bin"),
		scrapligooptions.WithBinTransportExtraArgs("-o ProxyCommand='foo' -P 1234"),
		scrapligooptions.WithBinTransportOverrideArgs("alll the args"),
		scrapligooptions.WithBinTransportSSHConfigFile("the/config/file"),
		scrapligooptions.WithBinTransportKnownHostsFile("the/known/hosts"),
		scrapligooptions.WithBinTransportStrictKey(),
		scrapligooptions.WithTermHeight(123),
		scrapligooptions.WithTermWidth(456),
	)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := d.GetOptions()
	if err != nil {
		t.Fatal(err)
	}

	testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", "transport_bin.json"))
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
