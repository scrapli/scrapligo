package testhelper

import (
	"flag"
	"strings"
	"testing"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
)

const all = "all"

var Functional = flag.Bool( //nolint:gochecknoglobals
	"functional", false, "perform functional tests")

var FunctionalPlatform = flag.String( //nolint:gochecknoglobals
	"platform", "all", "list comma sep platform(s) to target")

var FunctionalTransport = flag.String( //nolint:gochecknoglobals
	"transport", "all", "list comma sep transport(s) to target")

func RunPlatform(p string) bool {
	if *FunctionalPlatform == all {
		return true
	}

	platformTargetSplit := strings.Split(*FunctionalPlatform, ",")

	return util.StrInSlice(p, platformTargetSplit)
}

func RunTransport(t string) bool {
	if *FunctionalTransport == all {
		return true
	}

	platformTargetSplit := strings.Split(*FunctionalTransport, ",")

	return util.StrInSlice(t, platformTargetSplit)
}

type FunctionalTestHostConnData struct {
	Host       string
	Port       int
	TelnetPort int
}

func FunctionalTestHosts() map[string]*FunctionalTestHostConnData {
	return map[string]*FunctionalTestHostConnData{
		"cisco_iosxe": {
			Host:       "localhost",
			Port:       21022,
			TelnetPort: 21023,
		},
		"cisco_iosxr": {
			Host:       "localhost",
			Port:       23022,
			TelnetPort: 23023,
		},
		"cisco_nxos": {
			Host:       "localhost",
			Port:       22022,
			TelnetPort: 22023,
		},
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
