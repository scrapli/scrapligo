package netconf_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestAction(t *testing.T) {
	parentName := "action"

	cases := map[string]struct {
		description string
		action      string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple",
			action: `
			<system xmlns="urn:dummy:actions">
			  <reboot>
				<delay>5</delay>
			  </reboot>
			</system>`,
			options: []scrapligonetconf.Option{},
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

			r, err := n.Action(ctx, c.action, c.options...)
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
