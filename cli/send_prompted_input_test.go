package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
)

func TestSendPromptedInput(t *testing.T) {
	parentName := "send-prompted-input"

	cases := map[string]struct {
		description string
		postOpenF   func(t *testing.T, d *scrapligocli.Cli)
		input       string
		prompt      string
		response    string
		options     []scrapligocli.OperationOption
	}{
		"simple": {
			description: "simple input that requires no pagination",
			input:       "read -p \"Will you prompt me plz? \" answer",
			prompt:      "Will you prompt me plz?",
			response:    "nou",
			options: []scrapligocli.OperationOption{
				scrapligocli.WithRequestedMode("bash"),
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

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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
