package cfg_test

import (
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/cfg"
	"github.com/scrapli/scrapligo/util/cfgtesthelper"
)

func configSourceMap() map[string]string {
	return map[string]string{
		"cisco_iosxe":   "running",
		"cisco_iosxr":   "running",
		"cisco_nxos":    "running",
		"arista_eos":    "running",
		"juniper_junos": "running",
	}
}

func TestGetConfig(t *testing.T) {
	configSource := configSourceMap()

	for _, platform := range cfg.SupportedPlatforms() {
		f := cfgtesthelper.GetConfigTestHelper(t, platform, configSource[platform])
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}
