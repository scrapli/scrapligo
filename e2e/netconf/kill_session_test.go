package netconf_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligonetconf "github.com/scrapli/scrapligo/v2/netconf"
)

func TestKillSession(t *testing.T) {
	parentName := "kill-session"

	cases := map[string]struct {
		description string
		platform    string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple - kill the session",
			platform:    "nokia_srlinux",
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

				n2 := getNetconf(t, c.platform, transportName)

				_, err = n2.Open(ctx)
				if err != nil {
					t.Fatal(err)
				}

				s, err := n2.GetSessionID()
				if err != nil {
					t.Fatal(err)
				}

				r, err := n.KillSession(ctx, s)
				if err != nil {
					t.Fatal(err)
				}

				// ensure we free n2, if we dont could have segfaults if logging callbacks
				// end up poking something that got gc'd etc.
				n2ptr, ffiMapping := n2.GetPtr()
				ffiMapping.Shared.Free(n2ptr)

				assertResult(t, r, testGoldenPath)
			})
		}
	}
}
