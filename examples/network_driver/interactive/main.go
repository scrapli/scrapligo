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
		"ios-xe-mgmt.cisco.com",
		"cisco_iosxe",
		base.WithPort(8181),
		base.WithAuthStrictKey(false),
		base.WithAuthUsername("developer"),
		base.WithAuthPassword("C1sco12345"),
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
	defer d.Close()

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
}
