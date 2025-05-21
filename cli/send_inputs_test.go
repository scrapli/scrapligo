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
		postOpenF   func(t *testing.T, d *scrapligocli.Cli)
		inputs      []string
		options     []scrapligocli.OperationOption
	}{
		"simple-single-input": {
			description: "simple single input no pagination required",
			inputs:      []string{"show version | i Kern"},
			options:     []scrapligocli.OperationOption{},
		},
		"simple-multi-input": {
			description: "simple multi input no pagination required",
			inputs:      []string{"show version | i Kern", "show version | i Kern"},
			options:     []scrapligocli.OperationOption{},
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

			r, err := c.SendInputs(ctx, caseData.inputs, caseData.options...)
			if err != nil {
				t.Fatal(err)
			}

			assertResult(t, r, testGoldenPath)
		})
	}
}
