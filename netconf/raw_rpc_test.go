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

func TestRawRPC(t *testing.T) {
	parentName := "raw-rpc"

	cases := map[string]struct {
		description string
		payload     string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple",
			payload:     "<get-config><source><running/></source></get-config>",
		},
		"simple-extra-namespaces": {
			description: "simple",
			payload:     "<get-config><source><running/></source></get-config>",
			options: []scrapligonetconf.Option{
				scrapligonetconf.WithExtraNamespaces([][2]string{{"foo", "bar"}, {"baz", "qux"}}),
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

			r, err := n.RawRPC(ctx, c.payload, c.options...)
			if err != nil {
				t.Fatal(err)
			}

			assertResult(t, r, testGoldenPath)
		})
	}
}

func TestRawRPCCreateSubscription(t *testing.T) {
	parentName := "raw-rpc-create-subscription"

	cases := map[string]struct {
		description string
		payload     string
	}{
		"simple": {
			description: "simple",
			payload: `  <create-subscription xmlns="urn:ietf:params:xml:ns:netconf:notification:1.0">
    <stream>NETCONF</stream>
    <filter type="subtree">
      <counter-update xmlns="urn:boring:counter"/>
    </filter>
  </create-subscription>`,
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
