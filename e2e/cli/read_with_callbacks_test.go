package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
)

func TestReadWithCallbacks(t *testing.T) {
	parentName := "read-with-callbacks"

	cases := map[string]struct {
		description  string
		platform     string
		transports   []string
		initialInput string
		callbacks    []*scrapligocli.ReadCallback
	}{
		"simple-srl": {
			description:  "simple read with callbacks",
			platform:     scrapligocli.NokiaSrl.String(),
			transports:   []string{"bin", "ssh2"},
			initialInput: "show version | grep OS",
			callbacks: []*scrapligocli.ReadCallback{
				scrapligocli.NewReadCallback(
					"cb1",
					func(c *scrapligocli.Cli) error {
						return c.WriteAndReturn("show version | grep OS")
					},
					scrapligocli.WithContains("A:srl#"),
					scrapligocli.WithOnce(),
				),
				scrapligocli.NewReadCallback(
					"cb2",
					func(_ *scrapligocli.Cli) error {
						return nil
					},
					scrapligocli.WithContains("A:srl#"),
					scrapligocli.WithOnce(),
					scrapligocli.WithCompletes(),
				),
			},
		},
		"simple-eos": {
			description:  "simple read with callbacks",
			platform:     scrapligocli.AristaEos.String(),
			transports:   []string{"bin", "ssh2", "telnet"},
			initialInput: "show version | i Kernel",
			callbacks: []*scrapligocli.ReadCallback{
				scrapligocli.NewReadCallback(
					"cb1",
					func(c *scrapligocli.Cli) error {
						return c.WriteAndReturn("show version | i Kernel")
					},
					scrapligocli.WithContains("eos1#"),
					scrapligocli.WithOnce(),
				),
				scrapligocli.NewReadCallback(
					"cb2",
					func(_ *scrapligocli.Cli) error {
						return nil
					},
					scrapligocli.WithContains("eos1#"),
					scrapligocli.WithOnce(),
					scrapligocli.WithCompletes(),
				),
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

				r, err := c.ReadWithCallbacks(ctx, caseData.initialInput, caseData.callbacks...)
				if err != nil {
					t.Fatal(err)
				}

				assertResult(t, r, testGoldenPath)
			})
		}
	}
}
