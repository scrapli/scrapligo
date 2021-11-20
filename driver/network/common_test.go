package network_test

import (
	"testing"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
)

func platformCommandMapShort() map[string]string {
	return map[string]string{
		"cisco_iosxe":        "show run | i hostname",
		"cisco_iosxr":        "show run | i MgmtEth0",
		"cisco_nxos":         "show run | i scp-server",
		"arista_eos":         "show run | i ZTP",
		"juniper_junos":      "show configuration | match 10.0.0.15",
		"nokia_sros":         "show version",
		"nokia_sros_classic": "show version",
		"paloalto_panos":     "show clock",
	}
}

func platformCommandMapLong() map[string]string {
	return map[string]string{
		"cisco_iosxe":        "show run",
		"cisco_iosxr":        "show run",
		"cisco_nxos":         "show run",
		"arista_eos":         "show run",
		"juniper_junos":      "show configuration",
		"nokia_sros":         "show router interface",
		"nokia_sros_classic": "show router interface",
		"paloalto_panos":     "show templates",
	}
}

func platformConfigsMap() map[string][]string {
	return map[string][]string{
		"cisco_iosxe": {"interface loopback0", "description tacocat", "no interface loopback0"},
		"cisco_iosxr": {
			"interface loopback0",
			"description tacocat",
			"no interface loopback0",
			"commit",
		},
		"cisco_nxos": {"interface loopback0", "description tacocat", "no interface loopback0"},
		"arista_eos": {"interface loopback0", "description tacocat", "no interface loopback0"},
		"juniper_junos": {
			"set interfaces fxp0.0 description tacocat",
			"delete interfaces fxp0.0 description tacocat",
			"commit",
		},
		"nokia_sros": {
			`configure router interface "system" description "@ntdvps"`,
			"configure system",
			"location wide_internet",
			"commit",
		},
		"nokia_sros_classic": {
			`configure router interface "system" description "@ntdvps"`,
			"configure system",
			"location wide_internet",
		},
		"paloalto_panos": {
			"set display-name BLAH",
			"commit",
		},
	}
}

type functionalTestHostConnData struct {
	Host       string
	Port       int
	TelnetPort int
}

func functionalTestHosts() map[string]*functionalTestHostConnData {
	return map[string]*functionalTestHostConnData{
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

func newFunctionalTestDriver(
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
