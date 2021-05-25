package network_test

import (
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/driver/core"

	"github.com/scrapli/scrapligo/util/testhelper"
)

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
	}
}

func TestSendConfigs(t *testing.T) {
	commandsMap := platformConfigsMap()

	for _, platform := range core.SupportedPlatforms() {
		f := testhelper.SendConfigsTestHelper(platform, commandsMap[platform])
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}

func TestSendConfigsFromFile(t *testing.T) {
	for _, platform := range core.SupportedPlatforms() {
		f := testhelper.SendConfigsFromFileTestHelper(platform)
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}
