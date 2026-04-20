package options_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	scrapligocli "github.com/kentik/scrapligo/v2/cli"
	scrapligologging "github.com/kentik/scrapligo/v2/logging"
	scrapligooptions "github.com/kentik/scrapligo/v2/options"
	scrapligotesthelper "github.com/kentik/scrapligo/v2/testhelper"
)

func TestCommonOptions(t *testing.T) {
	d, err := scrapligocli.NewCli(
		"1.2.3.4",
		scrapligooptions.WithLoggerLevel(scrapligologging.Fatal),
		scrapligooptions.WithPort(1234),
	)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := d.GetOptions()
	if err != nil {
		t.Fatal(err)
	}

	testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", "common.json"))
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
