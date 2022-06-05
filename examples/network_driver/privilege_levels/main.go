package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/platform"
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

	d, err := p.GetNetworkDriver()
	if err != nil {
		fmt.Printf("failed to fetch network driver from the platform; error: %+v\n", err)
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

	// acquire configuration privilege level
	err = d.AcquirePriv("configuration")
	if err != nil {
		fmt.Printf("failed to acquire configuration privilege level; error: %+v\n", err)

		return
	}

	// fetch the prompt again to make sure we are in config mode
	prompt, err = d.GetPrompt()
	if err != nil {
		fmt.Printf("failed to get prompt; error: %+v\n", err)

		return
	}

	fmt.Printf("found prompt: %s\n\n\n", prompt)

	// now run a command that "should" be ran from privilege-exec mode -- you'll see that scrapligo
	// will automagically acquire privilege-exec -- this is the "default desired privilege level"
	// and this default priv level is *always* acquired before running "send commandX" methods.
	r, err := d.SendCommand("show run")
	if err != nil {
		fmt.Printf("failed to execute show command; error: %+v\n", err)

		return
	}

	fmt.Printf("got running config: %s\n\n\n", r.Result)
}
