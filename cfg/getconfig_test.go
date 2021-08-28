package cfg_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/scrapli/scrapligo/util/testhelper"

	"github.com/scrapli/scrapligo/cfg"
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

func testGetConfig(c *cfg.Cfg, source string, expected []byte) func(t *testing.T) {
	return func(t *testing.T) {
		r, cmdErr := c.GetConfig(source)
		if cmdErr != nil {
			t.Fatalf("failed running GetConfig: %v", cmdErr)
		}

		if r.Failed != nil {
			t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
		}

		if diff := cmp.Diff(r.Result, string(expected)); diff != "" {
			t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
		}
	}
}

func TestGetConfig(t *testing.T) {
	configSource := configSourceMap()

	for _, platform := range cfg.SupportedPlatforms() {
		sessionFile := fmt.Sprintf("../test_data/cfg/getconfig/%s", platform)
		expectedFile := fmt.Sprintf(
			"../test_data/cfg/getconfig/%s_expected",
			platform,
		)

		expected, expectedErr := os.ReadFile(expectedFile)
		if expectedErr != nil {
			t.Fatalf("failed opening expected output file '%s' err: %v", expectedFile, expectedErr)
		}

		d := testhelper.CreatePatchedDriver(t, sessionFile, platform)
		c := createCfgDriver(t, d, platform)

		f := testGetConfig(c, configSource[platform], expected)

		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}
