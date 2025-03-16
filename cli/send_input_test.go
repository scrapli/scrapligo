package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestSendInput(t *testing.T) {
	parentName := "send-input"

	cases := map[string]struct {
		description string
		postOpenF   func(t *testing.T, d *scrapligocli.Driver)
		input       string
		options     []scrapligocli.OperationOption
	}{
		"simple": {
			description: "simple input that requires no pagination",
			input:       "show version | i Kern",
			options:     []scrapligocli.OperationOption{},
		},
		"simple-requires-pagination": {
			description: "simple input that requires pagination",
			// dont include stuff like secretes as there are "$" that we will mistake for being a
			// prompt pattern because of test transport reading one byte at a time, so just show the
			// transceiver stuff since thats enough to require pagination!
			input:   "show running-config all | include snmp",
			options: []scrapligocli.OperationOption{},
		},
		"simple-already-in-non-default-mode": {
			description: "simple input executed in non-default mode we are already in",
			postOpenF: func(t *testing.T, d *scrapligocli.Driver) {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				_, err := d.EnterMode(ctx, "configuration")
				if err != nil {
					t.Fatal(err)
				}
			},
			input: "do show version | i Kern",
			options: []scrapligocli.OperationOption{
				scrapligocli.WithRequestedMode("configuration"),
			},
		},
		"simple-acquire-non-default-mode": {
			description: "simple input executed in freshly acquired non-default mode",
			input:       "do show version | i Kern",
			options: []scrapligocli.OperationOption{
				scrapligocli.WithRequestedMode("configuration"),
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

			d := getDriver(t, testFixturePath)
			defer closeDriver(t, d, testFixturePath)

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

				scrapligotesthelper.AssertEqual(t, 22, r.Port)
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
