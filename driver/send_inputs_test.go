package driver_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligodriver "github.com/scrapli/scrapligo/driver"
)

func TestSendInputs(t *testing.T) {
	parentName := "send-inputs"

	cases := map[string]struct {
		description string
		postOpenF   func(t *testing.T, d *scrapligodriver.Driver)
		inputs      []string
		options     []scrapligodriver.OperationOption
	}{
		"simple": {
			description: "simple input that requires no pagination",
			inputs:      []string{"show version | i Kern"},
			options:     []scrapligodriver.OperationOption{},
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

			_ = testGoldenPath
		})
	}
}
