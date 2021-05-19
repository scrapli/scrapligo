package network_test

import (
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/channel"

	"github.com/scrapli/scrapligo/driver/core"

	"github.com/scrapli/scrapligo/util/testhelper"
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

func TestSendInteractive(t *testing.T) {
	interactiveMap := platformInteractiveMap()

	for _, platform := range core.SupportedPlatforms() {
		if platform == "cisco_nxos" {
			// gotta figure out a interactive command on nxos
			continue
		}

		if platform == "juniper_junos" {
			// gotta figure out a interactive command on junos
			continue
		}

		f := testhelper.SendInteractiveTestHelper(platform, interactiveMap[platform])
		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}
