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

func TestSendInput(t *testing.T) {
	parentName := "send-input"

	cases := map[string]struct {
		description string
		postOpenF   func(t *testing.T, d *scrapligodriver.Driver)
		input       string
		options     []scrapligodriver.OperationOption
	}{
		"simple": {
			description: "simple input that requires no pagination",
			input:       "show version | i Kern",
			options:     []scrapligodriver.OperationOption{},
		},
		"simple-requires-pagination": {
			description: "simple input that requires pagination",
			input:       "show running-config all",
			options:     []scrapligodriver.OperationOption{},
		},
		"simple-non-default-mode": {
			description: "simple input executed in non-default mode",
			postOpenF: func(t *testing.T, d *scrapligodriver.Driver) {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				_, err := d.EnterMode(ctx, "configuration")
				if err != nil {
					t.Fatal(err)
				}
			},
			input:   "do show version | i Kern",
			options: []scrapligodriver.OperationOption{},
		},
		"simple-acquire-non-default-mode": {
			description: "simple input executed in freshly acquired non-default mode",
			input:       "do show version | i Kern",
			options: []scrapligodriver.OperationOption{
				scrapligodriver.WithRequestedMode("configuration"),
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

			r, err := d.SendInput(ctx, c.input)
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
