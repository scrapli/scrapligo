package options_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
	scrapligotesthelper "github.com/scrapli/scrapligo/v2/testhelper"
)

func TestCLIOptions(t *testing.T) {
	d, err := scrapligocli.NewCli(
		"1.2.3.4",
	)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := d.GetOptions()
	if err != nil {
		t.Fatal(err)
	}

	testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", "cli.json"))
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
