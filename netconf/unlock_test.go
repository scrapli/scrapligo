package netconf_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
)

func TestUnlock(t *testing.T) {
	parentName := "unlock"

	cases := map[string]struct {
		description string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple - lock the candidate config",
			options: []scrapligonetconf.Option{
				scrapligonetconf.WithDatastore(scrapligonetconf.DatastoreTypeCandidate),
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

			n := getNetconf(t, testFixturePath)

			_, err = n.Open(ctx)
			if err != nil {
				t.Fatal(err)
			}

			defer closeNetconf(t, n)

			_, err = n.Lock(ctx, c.options...)
			if err != nil {
				t.Fatal(err)
			}

			r, err := n.Unlock(ctx, c.options...)
			if err != nil {
				t.Fatal(err)
			}

			assertResult(t, r, testGoldenPath)
		})
	}
}
