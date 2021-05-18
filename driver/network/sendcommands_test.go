package network_test

import (
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/driver/core"

	"github.com/scrapli/scrapligo/util/testhelper"
)

func platformCommandsMap() map[string][]string {
	return map[string][]string{
		"cisco_iosxe":   {"show version", "show ip int brie"},
		"cisco_iosxr":   {"show version", "show ip int brie"},
		"cisco_nxos":    {"show version", "show ip int brie"},
		"arista_eos":    {"show version", "show ip int brie"},
		"juniper_junos": {"show version", "show interfaces terse"},
	}
}

func TestSendCommands(t *testing.T) {
	commandsMap := platformCommandsMap()

	for _, platform := range core.SupportedPlatforms() {
		f := testhelper.SendCommandsTestHelper(t, platform, commandsMap[platform])
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}

func TestSendCommandsFromFile(t *testing.T) {
	for _, platform := range core.SupportedPlatforms() {
		f := testhelper.SendCommandsFromFileTestHelper(t, platform)
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}
