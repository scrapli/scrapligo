package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/channel"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/driver/core"
)

func main() {
	// use the NewCoreDriver factory and pass in a platform argument
	d, err := core.NewCoreDriver(
		"localhost",
		"cisco_iosxe",
		base.WithPort(21022),
		base.WithAuthStrictKey(false),
		base.WithAuthUsername("vrnetlab"),
		base.WithAuthPassword("VR-netlab9"),
		base.WithAuthSecondary("VR-netlab9"),
	)

	if err != nil {
		fmt.Printf("failed to create driver; error: %+v\n", err)
		return
	}

	err = d.Open()
	if err != nil {
		fmt.Printf("failed to open driver; error: %+v\n", err)
		return
	}

	events := []*channel.SendInteractiveEvent{
		{
			ChannelInput:    "clear logging",
			ChannelResponse: "[confirm]",
			HideInput:       false,
		},
		{
			ChannelInput:    "",
			ChannelResponse: "",
			HideInput:       false,
		},
	}

	r, err := d.SendInteractive(events)
	if err != nil {
		fmt.Printf("failed to get prompt; error: %+v\n", err)
		return
	}
	fmt.Printf("interact response:\n%s\n", r.Result)

	err = d.Close()
	if err != nil {
		fmt.Printf("failed to close driver; error: %+v\n", err)
	}
}
