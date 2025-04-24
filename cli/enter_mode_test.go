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
		postOpenF     func(t *testing.T, d *scrapligocli.Driver)
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
			postOpenF: func(t *testing.T, d *scrapligocli.Driver) {
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
		"multi-stage-change-deescalate": {
			description: "enter mode with multi stage change 'deescalating' the mode'",
			postOpenF: func(t *testing.T, d *scrapligocli.Driver) {
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
