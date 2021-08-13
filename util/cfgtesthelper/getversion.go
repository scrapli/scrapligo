package cfgtesthelper

import (
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/cfg"
	"github.com/scrapli/scrapligo/util/testhelper"

	"github.com/scrapli/scrapligo/driver/core"
)

// GetVersionTestHelper helper function to handle get version tests.
func GetVersionTestHelper(t *testing.T, driverName, expectedVersion string) func(t *testing.T) {
	sessionFile := fmt.Sprintf("../test_data/cfg/getversion/%s", driverName)

	return func(t *testing.T) {
		d, driverErr := core.NewCoreDriver(
			"localhost",
			driverName,
			testhelper.WithPatchedTransport(sessionFile, t),
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

		r, cmdErr := c.GetVersion()
		if cmdErr != nil {
			t.Fatalf("failed running GetVersion: %v", cmdErr)
		}

		if r.Failed {
			t.Fatal("response object indicates failure")
		}

		if r.Result != expectedVersion {
			t.Errorf(
				"actual result and expected result do not match (-want +got):\n-%s +%s",
				r.Result,
				expectedVersion,
			)
		}
	}
}
