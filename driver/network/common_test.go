package network_test

import (
	"testing"

	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/util/testhelper"
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
	}
}

func createPatchedDriver(t *testing.T, sessionFile, platform string) *network.Driver {
	d, driverErr := core.NewCoreDriver(
		"localhost",
		platform,
		testhelper.WithPatchedTransport(sessionFile),
	)

	if driverErr != nil {
		t.Fatalf("failed creating test device: %v", driverErr)
	}

	return d
}
