package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
)

func customOnOpen(d *network.Driver) error {
	err := d.AcquirePriv(d.DefaultDesiredPriv)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("terminal length 0", nil)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("terminal width 512", nil)
	if err != nil {
		return err
	}

	return nil
}

func main() {
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

	// because of mostly copying python into go and being much less flexible w/ some magic typing
	// there is currently no way to pass an on open with the normal variadic args (because those
	// are *driver* options, not *network driver* options) -- probably this can be relaxed with
	// generics soon? in any case, you can still update/assign a custom on open function like so:
	d.OnOpen = customOnOpen

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
