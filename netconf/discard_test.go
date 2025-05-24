package netconf_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
)

func TestDiscard(t *testing.T) {
	parentName := "discard"

	cases := map[string]struct {
		description string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple - load a candidate config and discard it",
			options:     []scrapligonetconf.Option{},
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

			n := getNetconf(t, testFixturePath)

			_, err = n.Open(ctx)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_, _ = n.Close(ctx)
			}()

			r, err := n.Discard(ctx, c.options...)
			if err != nil {
				t.Fatal(err)
			}

			assertResult(t, r, testGoldenPath)
		})
	}
}
