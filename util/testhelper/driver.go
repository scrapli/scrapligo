package testhelper

import (
	"testing"

	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
)

func CreatePatchedDriver(t *testing.T, sessionFile, platform string) *network.Driver {
	d, driverErr := core.NewCoreDriver(
		"localhost",
		platform,
		WithPatchedTransport(sessionFile),
	)

	if driverErr != nil {
		t.Fatalf("failed creating test device: %v", driverErr)
	}

	return d
}
