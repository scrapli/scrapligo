package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/platform"

	"github.com/scrapli/scrapligo/channel"
)

func main() {
	p, err := platform.NewPlatform(
		// cisco_iosxe refers to the included cisco iosxe platform definition
		"cisco_iosxe",
		"sandbox-iosxe-latest-1.cisco.com",
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername("developer"),
		options.WithAuthPassword("C1sco12345"),
	)
	if err != nil {
		fmt.Printf("failed to create platform; error: %+v\n", err)

		return
	}

	// fetch the network driver instance from the platform. you need to call this method explicitly
	// because the platform may be generic or network -- by having the explicit method to fetch the
	// driver you can avoid having to type cast things yourself. if you had a generic driver based
	// platform you could call `GetGenericDriver` instead.
	d, err := p.GetNetworkDriver()
	if err != nil {
		fmt.Printf("failed to fetch network driver from the platform; error: %+v\n", err)

		return
	}

	err = d.Open()
	if err != nil {
		fmt.Printf("failed to open driver; error: %+v\n", err)

		return
	}

	defer d.Close()

	// fetch the prompt
	prompt, err := d.Channel.GetPrompt()
	if err != nil {
		fmt.Printf("failed to get prompt; error: %+v\n", err)

		return
	}

	fmt.Printf("found prompt: %s\n\n\n", prompt)

	// send some input
	output, err := d.Channel.SendInput("show version | i IOS")
	if err != nil {
		fmt.Printf("failed to send input to device; error: %+v\n", err)

		return
	}

	fmt.Printf("output received (SendInput):\n %s\n\n\n", output)

	// send an interactive input
	// SendInteractive expects a slice of `SendInteractiveEvent` objects
	events := make([]*channel.SendInteractiveEvent, 2)
	events[0] = &channel.SendInteractiveEvent{
		ChannelInput:    "clear logging",
		ChannelResponse: "[confirm]",
		HideInput:       false,
	}
	events[1] = &channel.SendInteractiveEvent{
		ChannelInput:    "",
		ChannelResponse: "#",
		HideInput:       false,
	}

	interactiveOutput, err := d.SendInteractive(events)
	if err != nil {
		fmt.Printf("failed to send interactive input to device; error: %+v\n", err)
	}

	fmt.Printf("output received (SendInteractive):\n %s\n\n\n", interactiveOutput.Result)

	// send a command -- as this is a driver created from a *platform* it will have some things
	// already done for us -- including disabling paging, so this command that would produce more
	// output than the default terminal lines will not cause any issues.
	r, err := d.SendCommand("show version")
	if err != nil {
		fmt.Printf("failed to send command; error: %+v\n", err)
		return
	}

	fmt.Printf(
		"sent command '%s', output received (SendCommand):\n %s\n\n\n",
		r.Input,
		r.Result,
	)
}
