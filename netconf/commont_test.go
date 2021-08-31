package netconf_test

import (
	"testing"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/netconf"
)

type functionalTestHostConnData struct {
	Host string
	Port int
}

func functionalTestHosts() map[string]*functionalTestHostConnData {
	return map[string]*functionalTestHostConnData{
		"cisco_iosxe_1_0": {
			Host: "localhost",
			Port: 21022,
		},
		"cisco_iosxe_1_1": {
			Host: "localhost",
			Port: 21830,
		},
		"cisco_iosxr_1_1": {
			Host: "localhost",
			Port: 23830,
		},
		"juniper_junos_1_0": {
			Host: "localhost",
			Port: 25022,
		},
	}
}

func newFunctionalTestDriver(
	t *testing.T,
	host, transportName string,
	port int,
) *netconf.Driver {
	d, driverErr := netconf.NewNetconfDriver(
		host,
		base.WithAuthUsername("boxen"),
		base.WithAuthPassword("b0x3N-b0x3N"),
		base.WithPort(port),
		base.WithTransportType(transportName),
		base.WithAuthStrictKey(false),
	)

	if driverErr != nil {
		t.Fatalf("failed creating test device: %v", driverErr)
	}

	return d
}
