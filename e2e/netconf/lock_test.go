package netconf_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
)

func TestLock(t *testing.T) {
	parentName := "lock"

	cases := map[string]struct {
		description string
		platform    string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple - lock the candidate config",
			platform:    "netopeer",
			options: []scrapligonetconf.Option{
				scrapligonetconf.WithDatastore(scrapligonetconf.DatastoreTypeCandidate),
			},
		},
	}

	for caseName, c := range cases {
		for _, transportName := range netconfTransports() {
			testName := fmt.Sprintf("%s-%s-%s-%s", parentName, caseName, c.platform, transportName)

			t.Run(testName, func(t *testing.T) {
				t.Logf("%s: starting", testName)

				testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", testName))
				if err != nil {
					t.Fatal(err)
				}

				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				n := getNetconf(t, c.platform, transportName)

				_, err = n.Open(ctx)
				if err != nil {
					t.Fatal(err)
				}

				defer func() {
					_, _ = n.Close(ctx)
				}()

				r, err := n.Lock(ctx, c.options...)
				if err != nil {
					t.Fatal(err)
				}

				_, err = n.Unlock(ctx, c.options...)
				if err != nil {
					t.Fatal(err)
				}

				assertResult(t, r, testGoldenPath)
			})
		}
	}
}
