package main

import (
	"fmt"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/options"
)

func main() {
	d, err := generic.NewDriver(
		"sandbox-iosxe-latest-1.cisco.com",
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername("developer"),
		options.WithAuthPassword("C1sco12345"),
		// the WithDefaultLogger option applies a simple log.Print logger at level info
		options.WithDefaultLogger(),
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

	r, err := d.SendCommand("show version | i Version")
	if err != nil {
		fmt.Printf("failed to run command; error: %+v\n", err)

		return
	}
	if r.Failed != nil {
		fmt.Printf("response object indicates failure: %+v\n", r.Failed)

		return
	}

	fmt.Printf("got some output: %s\n\n\n", r.Result)
}
