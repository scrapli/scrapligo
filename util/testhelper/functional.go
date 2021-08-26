package testhelper

import (
	"flag"
	"testing"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
)

var Functional = flag.Bool( //nolint:gochecknoglobals
	"functional", false, "perform functional tests")

type FunctionalTestHostConnData struct {
	Host       string
	Port       int
	TelnetPort int
}

func FunctionalTestHosts() map[string]*FunctionalTestHostConnData {
	return map[string]*FunctionalTestHostConnData{
		// "cisco_iosxe": {
		// 	Host:       "localhost",
		// 	Port:       21022,
		// 	TelnetPort: 21023,
		// },
		"arista_eos": {
			Host:       "localhost",
			Port:       24022,
			TelnetPort: 24023,
		},
		"juniper_junos": {
			Host:       "localhost",
			Port:       25022,
			TelnetPort: 25023,
		},
	}
}

func NewFunctionalTestDriver(
	t *testing.T,
	host, platform, transportName string,
	port int,
) *network.Driver {
	d, driverErr := core.NewCoreDriver(
		host,
		platform,
		base.WithAuthUsername("boxen"),
		base.WithAuthPassword("b0x3N-b0x3N"),
		base.WithAuthSecondary("b0x3N-b0x3N"),
		base.WithPort(port),
		base.WithTransportType(transportName),
		base.WithAuthStrictKey(false),
	)

	if driverErr != nil {
		t.Fatalf("failed creating test device: %v", driverErr)
	}

	return d
}
