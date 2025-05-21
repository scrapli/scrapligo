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
		content     string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple",
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

			r, err := n.EditData(ctx, c.content, c.options...)
			if err != nil {
				t.Fatal(err)
			}

			assertResult(t, r, testGoldenPath)
		})
	}
}
