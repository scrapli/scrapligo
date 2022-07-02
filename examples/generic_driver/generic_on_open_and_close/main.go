package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/options"
)

// custom "on X" functions (on open/close) accept a single argument of a pointer to a generic.Driver
// you can then access the driver to run any kind of setup/tear-down commands that you may need to
// prepare your device.
func customOnOpen(d *generic.Driver) error {
	_, err := d.SendCommand("terminal length 0")
	if err != nil {
		return err
	}

	_, err = d.SendCommand("terminal width 512")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	d, err := generic.NewDriver(
		"sandbox-iosxe-latest-1.cisco.com",
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername("developer"),
		options.WithAuthPassword("C1sco12345"),
		// apply your custom OnOpen function with the `WithOnOpen` option
		options.WithOnOpen(customOnOpen),
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

	prompt, err := d.GetPrompt()
	if err != nil {
		fmt.Printf("failed to get prompt; error: %+v\n", err)

		return
	}

	fmt.Printf("found prompt: %s\n\n\n", prompt)
}
