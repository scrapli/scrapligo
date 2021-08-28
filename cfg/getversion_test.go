package cfg_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util/testhelper"

	"github.com/scrapli/scrapligo/cfg"
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

func testGetVersion(c *cfg.Cfg, expected string) func(t *testing.T) {
	return func(t *testing.T) {
		r, cmdErr := c.GetVersion()
		if cmdErr != nil {
			t.Fatalf("failed running GetVersion: %v", cmdErr)
		}

		if r.Failed != nil {
			t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
		}

		if diff := cmp.Diff(r.Result, expected); diff != "" {
			t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
		}
	}
}

func TestGetVersion(t *testing.T) {
	versionMap := expectedVersionMap()

	for _, platform := range cfg.SupportedPlatforms() {
		sessionFile := fmt.Sprintf("../test_data/cfg/getversion/%s", platform)

		d := testhelper.CreatePatchedDriver(t, sessionFile, platform)
		c := createCfgDriver(t, d, platform)

		f := testGetVersion(c, versionMap[platform])

		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}
