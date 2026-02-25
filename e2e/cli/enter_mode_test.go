package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
)

func TestEnterMode(t *testing.T) {
	parentName := "enter-mode"

	cases := map[string]struct {
		description   string
		platform      string
		transports    []string
		postOpenF     func(t *testing.T, d *scrapligocli.Cli)
		requestedMode string
	}{
		"no-change-eos": {
			description:   "enter mode with no change required",
			platform:      scrapligocli.AristaEos.String(),
			transports:    []string{"bin", "ssh2", "telnet"},
			requestedMode: "privileged_exec",
		},
		"escalate-eos": {
			description:   "enter mode with single stage change 'escalating' the mode",
			platform:      scrapligocli.AristaEos.String(),
			transports:    []string{"bin", "ssh2", "telnet"},
			requestedMode: "configuration",
		},
		"multi-stage-change-escalate-eos": {
			description: "enter mode with multi stage change 'escalating' the mode'",
			platform:    scrapligocli.AristaEos.String(),
			transports:  []string{"bin", "ssh2", "telnet"},
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
			transports:  []string{"bin", "ssh2", "telnet"},
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
		"deescalate-eos": {
			description:   "enter mode with single stage change 'deescalating' the mode'",
			platform:      scrapligocli.AristaEos.String(),
			transports:    []string{"bin", "ssh2", "telnet"},
			requestedMode: "exec",
		},
		"no-change-srl": {
			description:   "enter mode with no change required",
			platform:      scrapligocli.NokiaSrlinux.String(),
			transports:    []string{"bin", "ssh2"},
			requestedMode: "exec",
		},
		"escalate-srl": {
			description:   "enter mode with single stage change 'escalating' the mode",
			platform:      scrapligocli.NokiaSrlinux.String(),
			transports:    []string{"bin", "ssh2"},
			requestedMode: "configuration",
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

				if caseData.postOpenF != nil {
					caseData.postOpenF(t, c)
				}

				r, err := c.EnterMode(ctx, caseData.requestedMode)
				if err != nil {
					t.Fatal(err)
				}

				assertResult(t, r, testGoldenPath)
			})
		}
	}
}
