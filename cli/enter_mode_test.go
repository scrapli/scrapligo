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
		postOpenF     func(t *testing.T, d *scrapligocli.Cli)
		requestedMode string
	}{
		"no-change": {
			description:   "enter mode with no change required",
			requestedMode: "privileged_exec",
		},
		"escalate": {
			description:   "enter mode with single stage change 'escalating' the mode",
			requestedMode: "configuration",
		},
		"deescalate": {
			description:   "enter mode with single stage change 'deescalating' the mode'",
			requestedMode: "exec",
		},
		"multi-stage-change-escalate": {
			description: "enter mode with multi stage change 'escalating' the mode'",
			postOpenF: func(t *testing.T, d *scrapligocli.Cli) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				_, err := d.EnterMode(ctx, "exec")
				if err != nil {
					t.Fatal(err)
				}
			},
			requestedMode: "configuration",
		},
		"multi-stage-change-deescalate": {
			description: "enter mode with multi stage change 'deescalating' the mode'",
			postOpenF: func(t *testing.T, d *scrapligocli.Cli) {
				t.Helper()

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				_, err := d.EnterMode(ctx, "configuration")
				if err != nil {
					t.Fatal(err)
				}
			},
			requestedMode: "exec",
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

			r, err := c.EnterMode(ctx, caseData.requestedMode)
			if err != nil {
				t.Fatal(err)
			}

			assertResult(t, r, testGoldenPath)
		})
	}
}
