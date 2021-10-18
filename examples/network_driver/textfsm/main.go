package main

import (
	"flag"
	"fmt"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
)

func main() {
	// File from https://github.com/networktocode/ntc-templates/blob/master/ntc_templates/templates/cisco_ios_show_version.textfsm
	arg := flag.String(
		"file",
		"examples/network_driver/textfsm/cisco_ios_show_version.textfsm",
		"argument from user",
	)
	flag.Parse()

	d, err := core.NewCoreDriver(
		"sandbox-iosxe-latest-1.cisco.com",
		"cisco_iosxe",
		base.WithAuthStrictKey(false),
		base.WithAuthUsername("developer"),
		base.WithAuthPassword("C1sco12345"),
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

	r, err := d.SendCommand("show version")
	if err != nil {
		fmt.Printf("failed to send command; error: %+v\n", err)
		return
	}

	parsedOut, err := r.TextFsmParse(*arg)
	if err != nil {
		fmt.Printf("failed to parse command; error: %+v\n", err)
		return
	}

	fmt.Printf("Hostname: %s\nSW Version: %s\nUptime: %s\n",
		parsedOut[0]["HOSTNAME"], parsedOut[0]["VERSION"], parsedOut[0]["UPTIME"])
}
