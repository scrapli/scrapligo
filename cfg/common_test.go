package cfg_test

import (
	"testing"

	"github.com/scrapli/scrapligo/cfg"
	"github.com/scrapli/scrapligo/driver/network"
)

func createCfgDriver(t *testing.T, d *network.Driver, platform string) *cfg.Cfg {
	openErr := d.Open()
	if openErr != nil {
		t.Fatalf("failed opening driver: %v", openErr)
	}

	c, cfgErr := cfg.NewCfgDriver(d, platform)

	if cfgErr != nil {
		t.Fatalf("failed creating cfg test device: %v", cfgErr)
	}

	prepareErr := c.Prepare()

	if prepareErr != nil {
		t.Fatalf("failed running prepare method: %v", prepareErr)
	}

	return c
}
