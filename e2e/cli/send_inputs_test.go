package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
)

func TestSendInputs(t *testing.T) {
	parentName := "send-inputs"

	cases := map[string]struct {
		description string
		platform    string
		transports  []string
		postOpenF   func(t *testing.T, d *scrapligocli.Cli)
		inputs      []string
		options     []scrapligocli.OperationOption
	}{
		"eos-single-input": {
			description: "simple input that requires no pagination",
			platform:    scrapligocli.AristaEos.String(),
			transports:  []string{"bin", "ssh2", "telnet"},
			inputs:      []string{"show version | i Kern"},
			options:     []scrapligocli.OperationOption{},
		},
		"eos-multi-input": {
			description: "simple input that requires pagination",
			platform:    scrapligocli.AristaEos.String(),
			transports:  []string{"bin", "ssh2", "telnet"},
			inputs:      []string{"show version | i Kern", "show run"},
			options:     []scrapligocli.OperationOption{},
		},
		"srl-multi-input": {
			description: "simple input that requires a mode change and requires pagination",
			platform:    scrapligocli.NokiaSrl.String(),
			transports:  []string{"bin", "ssh2"},
			inputs:      []string{"info system", "info"},
			options:     []scrapligocli.OperationOption{},
		},
	}

	for caseName, caseData := range cases {
		for _, transportName := range caseData.transports {
			testName := fmt.Sprintf("%s-%s-%s", parentName, caseName, transportName)

			t.Run(testName, func(t *testing.T) {
				t.Logf("%s: starting", testName)

				testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", testName))
				if err != nil {
					t.Fatal(err)
				}

				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				c := getCli(t, caseData.platform, transportName)
				defer func() {
					_, _ = c.Close(ctx)
				}()

				_, err = c.Open(ctx)
				if err != nil {
					t.Fatal(err)
				}

				if caseData.postOpenF != nil {
					caseData.postOpenF(t, c)
				}

				r, err := c.SendInputs(ctx, caseData.inputs, caseData.options...)
				if err != nil {
					t.Fatal(err)
				}

				assertResult(t, r, testGoldenPath)
			})
		}
	}
}
