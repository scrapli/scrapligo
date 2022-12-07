package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/options"
)

func main() {
	d, err := netconf.NewDriver(
		"sandbox-iosxe-latest-1.cisco.com",
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername("developer"),
		options.WithAuthPassword("C1sco12345"),
		options.WithPort(830),
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

	r, err := d.GetConfig("running")
	if err != nil {
		fmt.Printf("failed executing GetConfig; error: %+v\n", err)

		return
	}
	if r.Failed != nil {
		fmt.Printf("response object indicates failure: %+v\n", r.Failed)

		return
	}

	fmt.Printf("Config result: %s", r.Result)
}
