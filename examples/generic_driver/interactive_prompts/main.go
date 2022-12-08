package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/channel"
)

func main() {
	d, err := generic.NewDriver(
		"sandbox-iosxe-latest-1.cisco.com",
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername("developer"),
		options.WithAuthPassword("C1sco12345"),
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
	if r.Failed != nil {
		fmt.Printf("response object indicates failure: %+v\n", r.Failed)

		return
	}

	fmt.Printf("interact response:\n%s\n", r.Result)
}
