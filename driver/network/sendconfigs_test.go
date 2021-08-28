package network_test

import (
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/util/testhelper"

	"github.com/scrapli/scrapligo/driver/network"

	"github.com/scrapli/scrapligo/driver/core"
)

func testSendConfigs(d *network.Driver, config []string) func(t *testing.T) {
	return func(t *testing.T) {
		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening driver: %v", openErr)
		}

		r, cmdErr := d.SendConfigs(config)
		if cmdErr != nil {
			t.Fatalf("failed sending config: %v", cmdErr)
		}

		if r.Failed != nil {
			t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
		}
	}
}

func testSendConfigsFromFile(d *network.Driver, config string) func(t *testing.T) {
	return func(t *testing.T) {
		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening driver: %v", openErr)
		}

		r, cmdErr := d.SendConfigsFromFile(config)
		if cmdErr != nil {
			t.Fatalf("failed sending config: %v", cmdErr)
		}

		if r.Failed != nil {
			t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
		}
	}
}

func TestSendConfigs(t *testing.T) {
	configsMap := platformConfigsMap()

	for _, platform := range core.SupportedPlatforms() {
		sessionFile := fmt.Sprintf("../../test_data/driver/network/sendconfigs/%s", platform)

		d := testhelper.CreatePatchedDriver(t, sessionFile, platform)

		f := testSendConfigs(d, configsMap[platform])
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}

func TestSendConfigsFromFile(t *testing.T) {
	for _, platform := range core.SupportedPlatforms() {
		sessionFile := fmt.Sprintf("../../test_data/driver/network/sendconfigs/%s", platform)
		configs := fmt.Sprintf(
			"../../test_data/driver/network/sendconfigsfromfile/%s_configs",
			platform,
		)

		d := testhelper.CreatePatchedDriver(t, sessionFile, platform)

		f := testSendConfigsFromFile(d, configs)
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}
