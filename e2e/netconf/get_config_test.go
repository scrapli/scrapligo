package netconf_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestGetConfig(t *testing.T) {
	parentName := "get-config"

	cases := map[string]struct {
		description string
		platform    string
		options     []scrapligonetconf.Option
	}{
		"simple-eos": {
			description: "simple - get the running config",
			platform:    scrapligocli.AristaEos.String(),
		},
		"simple-srl": {
			description: "simple - get the running config",
			platform:    scrapligocli.NokiaSrl.String(),
		},
	}

	for caseName, c := range cases {
		for _, transportName := range getTransports() {
			if shouldSkip(c.platform, transportName) {
				continue
			}

			testName := fmt.Sprintf("%s-%s-%s", parentName, caseName, transportName)

			t.Run(testName, func(t *testing.T) {
				t.Logf("%s: starting", testName)

				testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", testName))
				if err != nil {
					t.Fatal(err)
				}

				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				n := getNetconf(t, c.platform, transportName)
				defer closeDriver(t, n)

				_, err = n.Open(ctx)
				if err != nil {
					t.Fatal(err)
				}

				r, err := n.GetConfig(ctx, c.options...)
				if err != nil {
					t.Fatal(err)
				}

				if *scrapligotesthelper.Update {
					scrapligotesthelper.WriteFile(
						t,
						testGoldenPath,
						scrapligotesthelper.CleanNetconfOutput(t, r.Result),
					)
				} else {
					cleanedActual := scrapligotesthelper.CleanNetconfOutput(t, r.Result)

					testGoldenContent := scrapligotesthelper.ReadFile(t, testGoldenPath)

					if !bytes.Equal(cleanedActual, testGoldenContent) {
						scrapligotesthelper.FailOutput(t, string(cleanedActual), testGoldenContent)
					}

					scrapligotesthelper.AssertNotDefault(t, r.StartTime)
					scrapligotesthelper.AssertNotDefault(t, r.EndTime)
					scrapligotesthelper.AssertNotDefault(t, r.ElapsedTimeSeconds)
					scrapligotesthelper.AssertNotDefault(t, r.Host)
					scrapligotesthelper.AssertNotDefault(t, r.ResultRaw)
				}
			})
		}
	}
}
