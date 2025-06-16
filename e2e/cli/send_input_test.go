package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
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
		transports  []string
		postOpenF   func(t *testing.T, d *scrapligocli.Cli)
		input       string
		options     []scrapligocli.Option
	}{
		"simple-eos": {
			description: "simple input that requires no pagination",
			platform:    scrapligocli.AristaEos.String(),
			transports:  []string{"bin", "ssh2", "telnet"},
			input:       "show version | i Ker",
			options:     []scrapligocli.Option{},
		},
		"simple-eos-pagination": {
			description: "simple input that requires pagination",
			platform:    scrapligocli.AristaEos.String(),
			transports:  []string{"bin", "ssh2", "telnet"},
			input:       "show run",
			options:     []scrapligocli.Option{},
		},
		"simple-eos-change-mode-and-pagination": {
			description: "simple input that requires a mode change and requires pagination",
			platform:    scrapligocli.AristaEos.String(),
			transports:  []string{"bin", "ssh2", "telnet"},
			input:       "show run",
			postOpenF: func(t *testing.T, d *scrapligocli.Cli) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				_, err := d.EnterMode(ctx, "configuration")
				if err != nil {
					t.Fatal(err)
				}
			},
			options: []scrapligocli.Option{
				scrapligocli.WithRequestedMode("privileged_exec"),
			},
		},
		"eos-retain-input": {
			description: "retain input in the final result",
			platform:    scrapligocli.AristaEos.String(),
			transports:  []string{"bin", "ssh2", "telnet"},
			input:       "show version | i Ker",
			options: []scrapligocli.Option{
				scrapligocli.WithRetainInput(),
			},
		},
		"eos-retain-trailing-prompt": {
			description: "retain trailing prompt in the final result",
			platform:    scrapligocli.AristaEos.String(),
			transports:  []string{"bin", "ssh2", "telnet"},
			input:       "show version | i Ker",
			options: []scrapligocli.Option{
				scrapligocli.WithRetainTrailingPrompt(),
			},
		},
		"eos-retain-trailing-all": {
			description: "retain trailing prompt in the final result",
			platform:    scrapligocli.AristaEos.String(),
			transports:  []string{"bin", "ssh2", "telnet"},
			input:       "show version | i Ker",
			options: []scrapligocli.Option{
				scrapligocli.WithRetainInput(),
				scrapligocli.WithRetainTrailingPrompt(),
			},
		},
		"simple-srl": {
			description: "simple input that requires no pagination",
			platform:    scrapligocli.NokiaSrl.String(),
			transports:  []string{"bin", "ssh2"},
			input:       "info interface mgmt0",
			options:     []scrapligocli.Option{},
		},
		"big-srl": {
			description: "simple input with a big output",
			platform:    scrapligocli.NokiaSrl.String(),
			transports:  []string{"bin", "ssh2"},
			input:       "info",
			options:     []scrapligocli.Option{},
		},
		"enormous-srl": {
			description: "simple input with an enormous output",
			platform:    scrapligocli.NokiaSrl.String(),
			input:       "info from state",
			options:     []scrapligocli.Option{},
		},
	}

	for caseName, caseData := range cases {
		for _, transportName := range caseData.transports {
			testName := fmt.Sprintf("%s-%s-%s", parentName, caseName, transportName)

			t.Run(testName, func(t *testing.T) {
				if *scrapligotesthelper.SkipSlow && slices.Contains(slowTests(), testName) {
					t.Skipf("skipping test %q due to skip slow flag", testName)
				}

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

				r, err := c.SendInput(ctx, caseData.input, caseData.options...)
				if err != nil {
					t.Fatal(err)
				}

				assertResult(t, r, testGoldenPath)
			})
		}
	}
}
