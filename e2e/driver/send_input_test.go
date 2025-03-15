package driver_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligodriver "github.com/scrapli/scrapligo/driver"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestSendInput(t *testing.T) {
	parentName := "send-input"

	cases := map[string]struct {
		description string
		platform    string
		transport   string
		postOpenF   func(t *testing.T, d *scrapligodriver.Driver)
		input       string
		options     []scrapligodriver.OperationOption
	}{
		"simple-srl": {
			description: "simple input that requires no pagination",
			platform:    scrapligodriver.NokiaSrl.String(),
			input:       "info interface mgmt0",
			options:     []scrapligodriver.OperationOption{},
		},
		"simple-eos": {
			description: "simple input that requires no pagination",
			platform:    scrapligodriver.AristaEos.String(),
			input:       "show version | i Kern",
			options:     []scrapligodriver.OperationOption{},
		},
	}

	for caseName, c := range cases {
		for _, transportName := range getTransports() {
			if shouldSkip(c.platform, c.transport) {
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

				d := getDriver(t, c.platform, transportName)
				defer closeDriver(t, d)

				_, err = d.Open(ctx)
				if err != nil {
					t.Fatal(err)
				}

				if c.postOpenF != nil {
					c.postOpenF(t, d)
				}

				r, err := d.SendInput(ctx, c.input, c.options...)
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
