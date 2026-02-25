package netconf_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligonetconf "github.com/scrapli/scrapligo/v2/netconf"
	scrapligotesthelper "github.com/scrapli/scrapligo/v2/testhelper"
)

func TestRawRPC(t *testing.T) {
	parentName := "raw-rpc"

	cases := map[string]struct {
		description string
		platform    string
		payload     string
		options     []scrapligonetconf.Option
	}{
		"simple": {
			description: "simple",
			platform:    "netopeer",
			payload:     "<get-config><source><running/></source></get-config>",
		},
		"simple-extra-namespaces": {
			description: "simple",
			platform:    "netopeer",
			payload:     "<get-config><source><running/></source></get-config>",
			options: []scrapligonetconf.Option{
				scrapligonetconf.WithExtraNamespaces([][2]string{{"foo", "bar"}, {"baz", "qux"}}),
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

				r, err := n.RawRPC(ctx, c.payload, c.options...)
				if err != nil {
					t.Fatal(err)
				}

				assertResult(t, r, testGoldenPath)
			})
		}
	}
}

func TestRawRPCCreateSubscription(t *testing.T) {
	parentName := "raw-rpc-create-subscription"

	cases := map[string]struct {
		description string
		platform    string
		payload     string
	}{
		"simple": {
			description: "simple",
			platform:    "netopeer",
			payload: `  <create-subscription xmlns="urn:ietf:params:xml:ns:netconf:notification:1.0">
    <stream>NETCONF</stream>
    <filter type="subtree">
      <counter-update xmlns="urn:boring:counter"/>
    </filter>
  </create-subscription>`,
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
}
