package cfgtesthelper

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/scrapli/scrapligo/cfg"
	"github.com/scrapli/scrapligo/util/testhelper"

	"github.com/scrapli/scrapligo/driver/core"
)

// GetConfigTestHelper helper function to handle get config tests.
func GetConfigTestHelper(t *testing.T, driverName, configSource string) func(t *testing.T) {
	sessionFile := fmt.Sprintf("../test_data/cfg/getconfig/%s", driverName)
	expectedFile := fmt.Sprintf(
		"../test_data/cfg/getconfig/%s_expected",
		driverName,
	)

	expected, expectedErr := os.ReadFile(expectedFile)
	if expectedErr != nil {
		t.Fatalf("failed opening expected output file '%s' err: %v", expectedFile, expectedErr)
	}

	return func(t *testing.T) {
		d, driverErr := core.NewCoreDriver(
			"localhost",
			driverName,
			testhelper.WithPatchedTransport(sessionFile),
		)

		if driverErr != nil {
			t.Fatalf("failed creating test device: %v", driverErr)
		}

		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening patched driver: %v", openErr)
		}

		c, cfgErr := cfg.NewCfgDriver(d, driverName)

		if cfgErr != nil {
			t.Fatalf("failed creating cfg test device: %v", driverErr)
		}

		prepareErr := c.Prepare()

		if prepareErr != nil {
			t.Fatalf("failed running prepare method: %v", driverErr)
		}

		r, cmdErr := c.GetConfig(configSource)
		if cmdErr != nil {
			t.Fatalf("failed running GetVersion: %v", cmdErr)
		}

		if r.Failed != nil {
			t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
		}

		// i have no idea where the null bit is getting read from... but it does? so we'll just remove
		// it for now...
		finalResult := string(bytes.Trim([]byte(r.Result), "\x00\x0a"))

		if diff := cmp.Diff(finalResult, string(expected)); diff != "" {
			t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
		}
	}
}
