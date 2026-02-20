package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
)

func TestGetPrompt(t *testing.T) {
	parentName := "get-prompt"

	cases := map[string]struct {
		description string
		platform    string
		transports  []string
	}{
		"simple-srl": {
			description: "simple get prompt",
			platform:    scrapligocli.NokiaSrlinux.String(),
			transports:  []string{"bin", "ssh2"},
		},
		"simple-eos": {
			description: "simple get prompt",
			platform:    scrapligocli.AristaEos.String(),
			transports:  []string{"bin", "ssh2", "telnet"},
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

				r, err := c.GetPrompt(ctx)
				if err != nil {
					t.Fatal(err)
				}

				assertResult(t, r, testGoldenPath)
			})
		}
	}
}
