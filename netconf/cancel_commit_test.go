package netconf_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestCancelCommit(t *testing.T) {
	parentName := "cancel-commit"

	cases := map[string]struct {
		description string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple",
			options:     []scrapligonetconf.Option{},
		},
	}

	for caseName, c := range cases {
		testName := fmt.Sprintf("%s-%s", parentName, caseName)

		t.Run(testName, func(t *testing.T) {
			t.Logf("%s: starting", testName)

			testFixturePath, err := filepath.Abs(fmt.Sprintf("./fixtures/%s", testName))
			if err != nil {
				t.Fatal(err)
			}

			testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", testName))
			if err != nil {
				t.Fatal(err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			n := getNetconf(t, testFixturePath)

			_, err = n.Open(ctx)
			if err != nil {
				t.Fatal(err)
			}

			defer closeNetconf(t, n)

			r, err := n.CancelCommit(ctx, c.options...)
			if err != nil {
				t.Fatal(err)
			}

			// rather than just use the handy assert result function we just do the stuff here
			// since this will be *failed* since there is no commit to cancel, *but* that does
			// show our rpc was/is valid/nicely formed, so thats what we care most about
			cleanedActual := scrapligotesthelper.CleanNetconfOutput(t, r.Result)

			// we can't just write the cleaned stuff to disk because then chunk sizes will be wrong if we
			// just do the lazy cleanup method we are doing (and cant stop wont stop)
			testGoldenContent := scrapligotesthelper.ReadFile(t, testGoldenPath)
			cleanedGolden := scrapligotesthelper.CleanNetconfOutput(t, string(testGoldenContent))

			if !bytes.Equal(cleanedActual, cleanedGolden) {
				scrapligotesthelper.FailOutput(t, cleanedActual, cleanedGolden)
			}

			scrapligotesthelper.AssertEqual(t, r.Port, 23830)
			scrapligotesthelper.AssertEqual(t, r.Host, testHost)
			scrapligotesthelper.AssertNotDefault(t, r.StartTime)
			scrapligotesthelper.AssertNotDefault(t, r.EndTime)
			scrapligotesthelper.AssertNotDefault(t, r.ElapsedTimeSeconds)
			scrapligotesthelper.AssertNotDefault(t, r.Host)
			scrapligotesthelper.AssertNotDefault(t, r.ResultRaw)
		})
	}
}
