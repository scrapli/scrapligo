package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
)

func TestSendPromptedInput(t *testing.T) {
	parentName := "send-prompted-input"

	cases := map[string]struct {
		description string
		platform    string
		transports  []string
		input       string
		prompt      string
		response    string
		options     []scrapligocli.Option
	}{
		"eos-simple": {
			description: "simple input that requires no pagination",
			platform:    scrapligocli.AristaEos.String(),
			transports:  []string{"bin", "ssh2", "telnet"},
			input:       "read -p \"Will you prompt me plz? \" answer",
			prompt:      "Will you prompt me plz?",
			response:    "nou",
			options: []scrapligocli.Option{
				scrapligocli.WithRequestedMode("bash"),
			},
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

				_, err = c.Open(ctx)
				if err != nil {
					t.Fatal(err)
				}

				defer func() {
					_, _ = c.Close(ctx)
				}()

				r, err := c.SendPromptedInput(
					ctx,
					caseData.input,
					caseData.prompt,
					caseData.response,
					caseData.options...)
				if err != nil {
					t.Fatal(err)
				}

				assertResult(t, r, testGoldenPath)
			})
		}
	}
}
