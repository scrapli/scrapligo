package network_test

import (
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/util/testhelper"

	"github.com/scrapli/scrapligo/driver/network"

	"github.com/scrapli/scrapligo/channel"

	"github.com/scrapli/scrapligo/driver/core"
)

func platformInteractiveMap() map[string][]*channel.SendInteractiveEvent {
	return map[string][]*channel.SendInteractiveEvent{
		"cisco_iosxe": {
			&channel.SendInteractiveEvent{
				ChannelInput:    "clear logging",
				ChannelResponse: "[confirm]",
				HideInput:       false,
			},
			&channel.SendInteractiveEvent{
				ChannelInput:    "",
				ChannelResponse: "",
				HideInput:       false,
			},
		},
		"arista_eos": {
			&channel.SendInteractiveEvent{
				ChannelInput:    "clear logging",
				ChannelResponse: "[confirm]",
				HideInput:       false,
			},
			&channel.SendInteractiveEvent{
				ChannelInput:    "",
				ChannelResponse: "",
				HideInput:       false,
			},
		},
		"cisco_iosxr": {
			{
				ChannelInput:    "clear logging",
				ChannelResponse: "Clear logging buffer",
				HideInput:       false,
			},
			{
				ChannelInput:    "y",
				ChannelResponse: "",
				HideInput:       false,
			},
		},
	}
}

func testSendInteractive(
	d *network.Driver,
	events []*channel.SendInteractiveEvent,
) func(t *testing.T) {
	return func(t *testing.T) {
		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening patched driver: %v", openErr)
		}

		r, interactErr := d.SendInteractive(events)
		if interactErr != nil {
			t.Fatalf("failed sending interactive: %v", interactErr)
		}

		if r.Failed != nil {
			t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
		}
	}
}

func TestSendInteractive(t *testing.T) {
	interactiveMap := platformInteractiveMap()

	for _, platform := range core.SupportedPlatforms() {
		if platform == "cisco_nxos" {
			// gotta figure out an interactive command on nxos
			continue
		}

		if platform == "juniper_junos" {
			// gotta figure out an interactive command on junos
			continue
		}

		if platform == "nokia_sros" {
			// gotta figure out an interactive command on sros
			continue
		}

		if platform == "nokia_sros_classic" {
			// gotta figure out an interactive command on sros
			continue
		}

		if platform == "paloalto_panos" {
			// gotta figure out an interactive command on sros
			continue
		}

		sessionFile := fmt.Sprintf("../../test_data/driver/network/sendinteractive/%s", platform)

		d := testhelper.CreatePatchedDriver(t, sessionFile, platform)

		f := testSendInteractive(d, interactiveMap[platform])
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}
