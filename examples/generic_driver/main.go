package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/driver/generic"
)

func main() {
	d, err := generic.NewGenericDriver(
		"localhost",
		base.WithPort(21022),
		base.WithAuthStrictKey(false),
		base.WithAuthUsername("vrnetlab"),
		base.WithAuthPassword("VR-netlab9"),
	)

	if err != nil {
		fmt.Printf("failed to create driver; error: %+v\n", err)
	}

	err = d.Open()
	if err != nil {
		fmt.Printf("failed to open driver; error: %+v\n", err)
	}

	// fetch the prompt
	prompt, err := d.GetPrompt()
	if err != nil {
		fmt.Printf("failed to get prompt; error: %+v\n", err)
	} else {
		fmt.Printf("found prompt: %s\n", prompt)
	}

	// send a command -- as this is "generic" there will have been no paging disabling, so have to
	// either disable paging yourself or send a command that will not make the device page the
	// output!
	r, err := d.SendCommand("show version | i uptime")
	if err != nil {
		fmt.Printf("failed to send command; error: %+v\n", err)
		return
	}
	fmt.Printf("sent command '%s', output received:\n %s\n", r.ChannelInput, r.Result)

	err = d.Close()
	if err != nil {
		fmt.Printf("failed to close driver; error: %+v\n", err)
	}
}
