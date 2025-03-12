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

func TestSendPromptedInput(t *testing.T) {
	parentName := "send-prompted-input"

	cases := map[string]struct {
		description string
		postOpenF   func(t *testing.T, d *scrapligodriver.Driver)
		input       string
		prompt      string
		response    string
		options     []scrapligodriver.OperationOption
	}{
		"simple": {
			description: "simple input that requires no pagination",
			input:       "read -p \"Will you prompt me plz? \" answer",
			prompt:      "Will you prompt me plz?",
			response:    "nou",
			options: []scrapligodriver.OperationOption{
				scrapligodriver.WithRequestedMode("bash"),
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
					[]byte(r.Result),
				)
			} else {
				testGoldenContent := scrapligotesthelper.ReadFile(t, testGoldenPath)

				if !bytes.Equal([]byte(r.Result), testGoldenContent) {
					t.Fatalf(
						"%s: actual and expected inputs do not match\nactual: %s\nexpected:%s",
						testName,
						r.Result,
						testGoldenContent,
					)
				}
			}
		})
	}
}
