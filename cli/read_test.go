package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
	scrapligotesthelper "github.com/scrapli/scrapligo/v2/testhelper"
)

func TestRead(t *testing.T) {
	parentName := "read"

	cases := map[string]struct {
		description string
		input       string
		options     []scrapligocli.Option
	}{
		"simple": {
			description: "read with standard size",
			input:       "show version | i Kern",
			options:     []scrapligocli.Option{},
		},
		"user-size": {
			description: "read with user provided size",
			input:       "show version | i Kern",
			options: []scrapligocli.Option{
				scrapligocli.WithReadSize(64),
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

			err = c.WriteAndReturn(caseData.input)
			if err != nil {
				t.Fatal(err)
			}

			time.Sleep(time.Second)

			b, err := c.Read(caseData.options...)
			if err != nil {
				t.Fatal(err)
			}

			if *scrapligotesthelper.Update {
				scrapligotesthelper.WriteFile(
					t,
					testGoldenPath,
					scrapligotesthelper.CleanCliOutput(t, string(b)),
				)

				return
			}

			testGoldenContent := scrapligotesthelper.ReadFile(t, testGoldenPath)

			if !bytes.Equal(b, testGoldenContent) {
				scrapligotesthelper.FailOutput(t, b, testGoldenContent)
			}
		})
	}
}
