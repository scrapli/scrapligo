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

func TestSendPromptedInput(t *testing.T) {
	parentName := "send-prompted-input"

	cases := map[string]struct {
		description string
		postOpenF   func(t *testing.T, d *scrapligocli.Driver)
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

			_, err = d.Open(ctx)
			if err != nil {
				t.Fatal(err)
			}

			defer closeDriver(t, d)

			if c.postOpenF != nil {
				c.postOpenF(t, d)
			}

			r, err := d.SendPromptedInput(ctx, c.input, c.prompt, c.response, c.options...)
			if err != nil {
				t.Fatal(err)
			}

			if *scrapligotesthelper.Update {
				scrapligotesthelper.WriteFile(
					t,
					testGoldenPath,
					scrapligotesthelper.CleanCliOutput(t, r.Result),
				)
			} else {
				assertResult(t, r, testGoldenPath)
			}
		})
	}
}
