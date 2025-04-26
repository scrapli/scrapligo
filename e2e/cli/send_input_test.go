package cli_test

import (
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
		platform    string
		postOpenF   func(t *testing.T, d *scrapligocli.Cli)
		input       string
		options     []scrapligocli.OperationOption
	}{
		"simple-srl": {
			description: "simple input that requires no pagination",
			platform:    scrapligocli.NokiaSrl.String(),
			input:       "info interface mgmt0",
			options:     []scrapligocli.OperationOption{},
		},
		"simple-eos": {
			description: "simple input that requires no pagination",
			platform:    scrapligocli.AristaEos.String(),
			input:       "show version | i Kern",
			options:     []scrapligocli.OperationOption{},
		},
		"big-srl": {
			description: "simple input with a big output",
			platform:    scrapligocli.NokiaSrl.String(),
			input:       "info",
			options:     []scrapligocli.OperationOption{},
		},
		// output file is literally 39MB, so... no, just no. but can be fun for testing!
		// if using need to set timeout > 140s or so (probably longer if in ci)
		// "enormous-srl": {
		// 	 description: "simple input with a big output",
		// 	 platform:    scrapligodriver.NokiaSrl.String(),
		// 	 input:       "info from state",
		// 	 options:     []scrapligodriver.OperationOption{},
		// },
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
						scrapligotesthelper.CleanCliOutput(t, r.Result()),
					)
				} else {
					assertResult(t, r, testGoldenPath)
				}
			})
		}
	}
}
