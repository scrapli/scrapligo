package driver_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	scrapligodriver "github.com/scrapli/scrapligo/driver"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestMain(m *testing.M) {
	scrapligotesthelper.Flags()

	os.Exit(m.Run())
}

func getDriver(t *testing.T, f string) *scrapligodriver.Driver {
	opts := []scrapligodriver.Option{
		scrapligodriver.WithUsername("admin"),
		scrapligodriver.WithPassword("admin"),
		scrapligodriver.WithLookupKeyValue("enable", "libscrapli"),
	}

	if *scrapligotesthelper.Record {
		opts = append(
			opts,
			scrapligodriver.WithPort(22022),
			scrapligodriver.WithSessionRecorderPath(f),
		)
	} else {
		opts = append(
			opts,
			scrapligodriver.WithTransportKind(scrapligodriver.TransportKindTest),
			scrapligodriver.WithTestTransportF(f),
			scrapligodriver.WithReadSize(1),
		)
	}

	d, err := scrapligodriver.NewDriver(
		string(scrapligodriver.AristaEos),
		"localhost",
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func TestEnterMode(t *testing.T) {
	parentName := "enter-mode"

	cases := map[string]struct {
		description   string
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

			r, err := d.EnterMode(ctx, c.requestedMode)
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
