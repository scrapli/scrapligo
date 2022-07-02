package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/channel"
)

func main() {
	d, err := network.NewDriver(
		"sandbox-iosxe-latest-1.cisco.com",
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername("developer"),
		options.WithAuthPassword("C1sco12345"),
		// network drivers *must* have a default desired privilege level and privileges set!
		options.WithDefaultDesiredPriv("privilege-exec"),
		options.WithPrivilegeLevels(map[string]*network.PrivilegeLevel{
			"exec": {
				Pattern:        `(?im)^[\w.\-@/:]{1,63}>$`,
				Name:           "exec",
				PreviousPriv:   "",
				Deescalate:     "",
				Escalate:       "",
				EscalateAuth:   false,
				EscalatePrompt: "",
			},
			"privilege-exec": {
				Pattern:        `(?im)^[\w.\-@/:]{1,63}#$`,
				Name:           "privilege-exec",
				PreviousPriv:   "exec",
				Deescalate:     "disable",
				Escalate:       "enable",
				EscalateAuth:   true,
				EscalatePrompt: `(?im)^(?:enable\s){0,1}password:\s?$`,
			},
			"configuration": {
				Pattern:        `(?im)^[\w.\-@/:]{1,63}\([\w.\-@/:+]{0,32}\)#$`,
				NotContains:    []string{"tcl)"},
				Name:           "configuration",
				PreviousPriv:   "privilege-exec",
				Deescalate:     "end",
				Escalate:       "configure terminal",
				EscalateAuth:   false,
				EscalatePrompt: "",
			},
		}),
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

		return
	}

	fmt.Printf("output received (SendInteractive):\n %s\n\n\n", interactiveOutput.Result)

	// send a command -- as this is a "base" driver (meaning there is no context of the type of
	// device we are connecting to) there will have been no paging disabling, so have to
	// either disable paging yourself or send a command that will not make the device page the
	// output!
	r, err := d.SendCommand("show version | i uptime")
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
