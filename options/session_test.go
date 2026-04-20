package options_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
	scrapligooptions "github.com/scrapli/scrapligo/v2/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/v2/testhelper"
)

func TestSessionOptions(t *testing.T) {
	d, err := scrapligocli.NewCli(
		"1.2.3.4",
		scrapligooptions.WithReadSize(999),
		scrapligooptions.WithReadMinDelay(123),
		scrapligooptions.WithReadMaxDelay(456),
		scrapligooptions.WithReturnChar("return"),
		scrapligooptions.WithOperationTimeout(time.Hour),
		scrapligooptions.WithOperationMaxSearchDepth(9876),
		scrapligooptions.WithSessionRecorderPath("recorderoutput"),
	)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := d.GetOptions()
	if err != nil {
		t.Fatal(err)
	}

	testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", "session.json"))
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
