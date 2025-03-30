package netconf_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestRawRPC(t *testing.T) {
	parentName := "raw-rpc"

	cases := map[string]struct {
		description string
		payload     string
	}{
		"simple": {
			description: "simple",
			payload:     "<get-config><source><running/></source></get-config>",
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

			n := getDriver(t, testFixturePath)

			_, err = n.Open(ctx)
			if err != nil {
				t.Fatal(err)
			}

			defer closeDriver(t, n)

			r, err := n.RawRPC(ctx, c.payload)
			if err != nil {
				t.Fatal(err)
			}

			if *scrapligotesthelper.Update {
				scrapligotesthelper.WriteFile(
					t,
					testGoldenPath,
					scrapligotesthelper.CleanNetconfOutput(t, r.Result),
				)
			} else {
				assertResult(t, r, testGoldenPath)
			}
		})
	}
}
