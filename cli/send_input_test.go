package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
)

func TestSendInput(t *testing.T) {
	parentName := "send-input"

	cases := map[string]struct {
		description string
		postOpenF   func(t *testing.T, d *scrapligocli.Cli)
		input       string
		options     []scrapligocli.Option
	}{
		"simple": {
			description: "simple input that requires no pagination",
			input:       "show version | i Kern",
			options:     []scrapligocli.Option{},
		},
		"simple-requires-pagination": {
			description: "simple input that requires pagination",
			// dont include stuff like secretes as there are "$" that we will mistake for being a
			// prompt pattern because of test transport reading one byte at a time, so just show the
			// transceiver stuff since thats enough to require pagination!
			input:   "show running-config all | include snmp",
			options: []scrapligocli.Option{},
		},
		"simple-already-in-non-default-mode": {
			description: "simple input executed in non-default mode we are already in",
			postOpenF: func(t *testing.T, d *scrapligocli.Cli) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				_, err := d.EnterMode(ctx, "configuration")
				if err != nil {
					t.Fatal(err)
				}
			},
			input: "do show version | i Kern",
			options: []scrapligocli.Option{
				scrapligocli.WithRequestedMode("configuration"),
			},
		},
		"simple-acquire-non-default-mode": {
			description: "simple input executed in freshly acquired non-default mode",
			input:       "do show version | i Kern",
			options: []scrapligocli.Option{
				scrapligocli.WithRequestedMode("configuration"),
			},
		},
		"simple-input-handling-exact": {
			description: "simple with exact input handling mode",
			input:       "show version | i Kern",
			options: []scrapligocli.Option{
				scrapligocli.WithInputHandling(scrapligocli.InputHandlingExact),
			},
		},
		"simple-input-handling-ignore": {
			description: "simple with ignore input handling mode",
			input:       "show version | i Kern",
			options: []scrapligocli.Option{
				scrapligocli.WithInputHandling(scrapligocli.InputHandlingIgnore),
			},
		},
		"simple-retain-input": {
			description: "simple with retain input",
			input:       "show version | i Kern",
			options: []scrapligocli.Option{
				scrapligocli.WithRetainInput(),
			},
		},
		"simple-retain-trailing-prompt": {
			description: "simple with retain trailing prompt",
			input:       "show version | i Kern",
			options: []scrapligocli.Option{
				scrapligocli.WithRetainTrailingPrompt(),
			},
		},
		"simple-retain-all": {
			description: "simple with retain input and trailing prompt",
			input:       "show version | i Kern",
			options: []scrapligocli.Option{
				scrapligocli.WithRetainInput(),
				scrapligocli.WithRetainTrailingPrompt(),
			},
		},
	}

	for caseName, caseData := range cases {
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

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			c := getCli(t, testFixturePath)

			_, err = c.Open(ctx)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_, _ = c.Close(ctx)
			}()

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
