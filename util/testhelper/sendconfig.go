package testhelper

import (
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/driver/core"
)

// SendConfigTestHelper helper function to handle send config tests.
func SendConfigTestHelper(driverName, config string) func(t *testing.T) {
	sessionFile := fmt.Sprintf("../../test_data/driver/network/sendconfigs/%s", driverName)

	return func(t *testing.T) {
		d, driverErr := core.NewCoreDriver(
			"localhost",
			driverName,
			WithPatchedTransport(sessionFile),
		)

		if driverErr != nil {
			t.Fatalf("failed creating test device: %v", driverErr)
		}

		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening patched driver: %v", openErr)
		}

		r, cmdErr := d.SendConfig(config)
		if cmdErr != nil {
			t.Fatalf("failed sending configs: %v", cmdErr)
		}

		if r.Failed {
			t.Fatal("response object indicates failure")
		}
	}
}
