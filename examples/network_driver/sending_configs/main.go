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

	// sending configs is just as easy as sending commands, scrapligo will auto acquire the
	// "configuration" privilege level, if your platform does *not* have a "configuration" priv
	// level things will fail! if there is an alternative to configuration (such as "exclusive"),
	// you can explicitly execute configs in that privilege level with the
	// `opoptions.WithPrivilegeLevel` option.
	r, err := d.SendConfigs([]string{"interface loopback999", "description tacocat"})
	if err != nil {
		fmt.Printf("failed to open driver; error: %+v\n", err)

		return
	}
	if r.Failed != nil {
		fmt.Printf("response object indicates failure: %+v\n", r.Failed)

		return
	}

	fmt.Printf("sending configs took %f seconds", r.ElapsedTime)
}
