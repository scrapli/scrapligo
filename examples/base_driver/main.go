package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/driver/base"
)

func main() {
	d, err := base.NewDriver(
		"sandbox-iosxe-latest-1.cisco.com",
		base.WithAuthStrictKey(false),
		base.WithAuthUsername("developer"),
		base.WithAuthPassword("C1sco12345"),
	)

	if err != nil {
		fmt.Printf("failed to create driver; error: %+v\n", err)
	}

	err = d.Open()
	if err != nil {
		fmt.Printf("failed to open driver; error: %+v\n", err)
	}
	defer d.Close()

	// fetch the prompt
	prompt, err := d.Channel.GetPrompt()
	if err != nil {
		fmt.Printf("failed to get prompt; error: %+v\n", err)
	} else {
		fmt.Printf("found prompt: %s\n\n\n", prompt)
	}

	// send some input
	// at the "base" level there are no convenience wrappers around the channel supporting options,
	// so you have to specify all the parameters when using the channel directly
	output, err := d.Channel.SendInput("show version | i IOS", true, false, -1)
	if err != nil {
		fmt.Printf("failed to send input to device; error: %+v\n", err)
	} else {
		fmt.Printf("output received (SendInput):\n %s\n\n\n", output)
	}

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

	// you can access the channel directly, however there are no convenience wrappers around the
	// channel supporting options, so you have to specify all the parameters when using it directly
	interactiveOutput, err := d.Channel.SendInteractive(events, nil, -1)
	if err != nil {
		fmt.Printf("failed to send interactive input to device; error: %+v\n", err)
	} else {
		fmt.Printf("output received (SendInteractive):\n %s\n\n\n", interactiveOutput)
	}

	// send a command -- as this is "base" there will have been no paging disabling, so have to
	// either disable paging yourself or send a command that will not make the device page the
	// output!
	r, err := d.SendCommand("show version | i uptime")
	if err != nil {
		fmt.Printf("failed to send command; error: %+v\n", err)
		return
	}
	fmt.Printf(
		"sent command '%s', output received (SendCommand):\n %s\n\n\n",
		r.ChannelInput,
		r.Result,
	)
}
