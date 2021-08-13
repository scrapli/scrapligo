package cfg_test

import (
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/cfg"
	"github.com/scrapli/scrapligo/util/cfgtesthelper"
)

func expectedVersionMap() map[string]string {
	return map[string]string{
		"cisco_iosxe":   "16.12.03",
		"cisco_iosxr":   "6.5.3",
		"cisco_nxos":    "9.2(4)",
		"arista_eos":    "4.22.1F",
		"juniper_junos": "17.3R2.10",
	}
}

func TestGetVersion(t *testing.T) {
	versionMap := expectedVersionMap()

	for _, platform := range cfg.SupportedPlatforms() {
		f := cfgtesthelper.GetVersionTestHelper(t, platform, versionMap[platform])
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}
