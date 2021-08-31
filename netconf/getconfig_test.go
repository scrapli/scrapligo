package netconf_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/scrapli/scrapligo/netconf"

	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util/testhelper"
)

func testFunctionalGetConfig(d *netconf.Driver) func(t *testing.T) {
	return func(t *testing.T) {
		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening driver: %v", openErr)
		}

		r, cmdErr := d.GetConfig("running")
		if cmdErr != nil {
			t.Fatalf("failed sending config: %v", cmdErr)
		}

		if r.Failed != nil {
			t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
		}
	}
}

func TestFunctionalGetConfig(t *testing.T) {
	if !*testhelper.Functional {
		t.Skip("skip: functional tests skipped unless the '-functional' flag is passed")
	}

	testHosts := functionalTestHosts()

	for _, transportName := range transport.SupportedNetconfTransports() {
		if !testhelper.RunTransport(transportName) {
			t.Logf("skip; transport %s deselected for testing\n", transportName)
			continue
		}

		// for now just making sure the damn thing runs eventually load up expected output and
		// compare actual<>expected like the other tests.
		for platform, connectionData := range testHosts {
			d := newFunctionalTestDriver(t, connectionData.Host, transportName, connectionData.Port)

			if strings.Contains(platform, "junos") {
				d.NetconfChannel.ForceSelfClosingTag = true
			}

			f := testFunctionalGetConfig(d)

			t.Run(fmt.Sprintf("Platform=%s;Transport=%s", platform, transportName), f)
		}
	}
}
