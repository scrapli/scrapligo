package testhelper

import (
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/channel"

	"github.com/scrapli/scrapligo/driver/core"
)

// SendInteractiveTestHelper helper function to handle send interactive event tests.
func SendInteractiveTestHelper(
	driverName string,
	interactEvents []*channel.SendInteractiveEvent,
) func(t *testing.T) {
	sessionFile := fmt.Sprintf("../../test_data/driver/network/sendinteractive/%s", driverName)

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

		r, interactErr := d.SendInteractive(interactEvents)
		if interactErr != nil {
			t.Fatalf("failed sending interactive: %v", interactErr)
		}

		if r.Failed != nil {
			t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
		}
	}
}
