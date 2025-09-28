package netconf_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
)

func TestKillSession(t *testing.T) {
	parentName := "kill-session"

	cases := map[string]struct {
		description string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple - kill the session",
		},
	}

	for caseName := range cases {
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

			n := getNetconfSrl(t, testFixturePath)

			_, err = n.Open(ctx)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_, _ = n.Close(ctx)
			}()

			n2 := getNetconfSrl(t, testFixturePath)

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
