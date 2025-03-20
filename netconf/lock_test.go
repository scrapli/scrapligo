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

func TestLock(t *testing.T) {
	parentName := "lock"

	cases := map[string]struct {
		description string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple - lock the candidate config",
			options: []scrapligonetconf.Option{
				scrapligonetconf.WithDatastore(scrapligonetconf.DatastoreTypeCandidate),
			},
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
			defer closeNetconf(t, n, testFixturePath)

			_, err = n.Open(ctx)
			if err != nil {
				t.Fatal(err)
			}

			r, err := n.Lock(ctx, c.options...)
			if err != nil {
				t.Fatal(err)
			}

			_, err = n.Unlock(ctx, c.options...)
			if err != nil {
				t.Fatal(err)
			}

			if *scrapligotesthelper.Update {
				scrapligotesthelper.WriteFile(
					t,
					testGoldenPath,
					[]byte(r.Result),
				)
			} else {
				testGoldenContent := scrapligotesthelper.ReadFile(t, testGoldenPath)

				if !bytes.Equal([]byte(r.Result), testGoldenContent) {
					scrapligotesthelper.FailOutput(t, r.Result, testGoldenContent)
				}

				scrapligotesthelper.AssertEqual(t, 22830, r.Port)
				scrapligotesthelper.AssertEqual(t, testHost, r.Host)
				scrapligotesthelper.AssertNotDefault(t, r.StartTime)
				scrapligotesthelper.AssertNotDefault(t, r.EndTime)
				scrapligotesthelper.AssertNotDefault(t, r.ElapsedTimeSeconds)
				scrapligotesthelper.AssertNotDefault(t, r.Host)
				scrapligotesthelper.AssertNotDefault(t, r.ResultRaw)
			}
		})
	}
}
