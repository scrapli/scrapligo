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

func TestEnterMode(t *testing.T) {
	parentName := "enter-mode"

	cases := map[string]struct {
		description   string
		platform      string
		postOpenF     func(t *testing.T, d *scrapligocli.Cli)
		requestedMode string
	}{
		"no-change-eos": {
			description:   "enter mode with no change required",
			platform:      scrapligocli.AristaEos.String(),
			requestedMode: "privileged_exec",
		},
		"escalate-eos": {
			description:   "enter mode with single stage change 'escalating' the mode",
			platform:      scrapligocli.AristaEos.String(),
			requestedMode: "configuration",
		},
		"deescalate-eos": {
			description:   "enter mode with single stage change 'deescalating' the mode'",
			platform:      scrapligocli.AristaEos.String(),
			requestedMode: "exec",
		},
		"multi-stage-change-escalate-eos": {
			description: "enter mode with multi stage change 'escalating' the mode'",
			platform:    scrapligocli.AristaEos.String(),
			postOpenF: func(t *testing.T, d *scrapligocli.Cli) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				_, err := d.EnterMode(ctx, "exec")
				if err != nil {
					t.Fatal(err)
				}
			},
			requestedMode: "configuration",
		},
		"multi-stage-change-deescalate-eos": {
			description: "enter mode with multi stage change 'deescalating' the mode'",
			platform:    scrapligocli.AristaEos.String(),
			postOpenF: func(t *testing.T, d *scrapligocli.Cli) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				_, err := d.EnterMode(ctx, "configuration")
				if err != nil {
					t.Fatal(err)
				}
			},
			requestedMode: "exec",
		},
	}

	for caseName, c := range cases {
		for _, transportName := range getTransports() {
			if shouldSkip(c.platform, transportName) {
				continue
			}

			testName := fmt.Sprintf("%s-%s-%s", parentName, caseName, transportName)

			t.Run(testName, func(t *testing.T) {
				t.Logf("%s: starting", testName)

				testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", testName))
				if err != nil {
					t.Fatal(err)
				}

				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				d := getCli(t, c.platform, transportName)
				defer closeCli(t, d)

				_, err = d.Open(ctx)
				if err != nil {
					t.Fatal(err)
				}

				if c.postOpenF != nil {
					c.postOpenF(t, d)
				}

				r, err := d.EnterMode(ctx, c.requestedMode)
				if err != nil {
					t.Fatal(err)
				}

				if *scrapligotesthelper.Update {
					scrapligotesthelper.WriteFile(
						t,
						testGoldenPath,
						scrapligotesthelper.CleanCliOutput(t, r.Result()),
					)
				} else {
					assertResult(t, r, testGoldenPath)
				}
			})
		}
	}
}
