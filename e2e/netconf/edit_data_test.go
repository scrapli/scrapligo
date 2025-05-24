package netconf_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
)

func TestEditData(t *testing.T) {
	parentName := "edit-data"

	cases := map[string]struct {
		description string
		platform    string
		content     string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple",
			platform:    "netopeer",
			content: `<system xmlns="urn:some:data">
        <hostname>my-router</hostname>
        <interfaces>
          <name>eth0</name>
          <enabled>true</enabled>
        </interfaces>
        <interfaces>
          <name>eth1</name>
          <enabled>false</enabled>
        </interfaces>
      </system>`,
			options: []scrapligonetconf.Option{},
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

				r, err := n.EditData(ctx, c.content, c.options...)
				if err != nil {
					t.Fatal(err)
				}

				assertResult(t, r, testGoldenPath)
			})
		}
	}
}
