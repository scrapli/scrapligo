package network_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/scrapli/scrapligo/util/testhelper"

	"github.com/scrapli/scrapligo/driver/network"

	"github.com/scrapli/scrapligo/driver/core"
)

func testSendConfig(d *network.Driver, config string) func(t *testing.T) {
	return func(t *testing.T) {
		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening driver: %v", openErr)
		}

		r, cmdErr := d.SendConfig(config)
		if cmdErr != nil {
			t.Fatalf("failed sending config: %v", cmdErr)
		}

		if r.Failed != nil {
			t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
		}
	}
}

func TestSendConfig(t *testing.T) {
	configsMap := platformConfigsMap()

	for _, platform := range core.SupportedPlatforms() {
		sessionFile := fmt.Sprintf("../../test_data/driver/network/sendconfigs/%s", platform)

		d := testhelper.CreatePatchedDriver(t, sessionFile, platform)

		f := testSendConfig(d, strings.Join(configsMap[platform], "\n"))
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}
