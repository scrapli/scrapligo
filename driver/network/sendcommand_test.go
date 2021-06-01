package network_test

import (
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/driver/core"

	"github.com/scrapli/scrapligo/util/testhelper"
)

func platformCommandMap() map[string]string {
	return map[string]string{
		"cisco_iosxe":        "show version",
		"cisco_iosxr":        "show version",
		"cisco_nxos":         "show version",
		"arista_eos":         "show version",
		"juniper_junos":      "show version",
		"nokia_sros":         "show version",
		"nokia_sros_classic": "show version",
	}
}

func TestSendCommand(t *testing.T) {
	commandMap := platformCommandMap()

	for _, platform := range core.SupportedPlatforms() {
		f := testhelper.SendCommandTestHelper(t, platform, commandMap[platform])
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}
