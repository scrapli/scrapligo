package main

import (
	"flag"
	"fmt"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
)

// const commandsFile = "commandsfile"

func main() {
	arg := flag.String("file", "examples/network_driver/simple/commandsfile", "argument from user")
	flag.Parse()

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

	err = d.Open()
	if err != nil {
		fmt.Printf("failed to open driver; error: %+v\n", err)
		return
	}
	defer d.Close()

	// fetch the prompt
	prompt, err := d.GetPrompt()
	if err != nil {
		fmt.Printf("failed to get prompt; error: %+v\n", err)
	} else {
		fmt.Printf("found prompt: %s\n\n\n", prompt)
	}

	// send some commands from a file
	mr, err := d.SendCommandsFromFile(*arg)
	if err != nil {
		fmt.Printf("failed to send commands from file; error: %+v\n", err)
		return
	}
	for _, r := range mr.Responses {
		fmt.Printf("sent command '%s', output received:\n %s\n\n\n", r.ChannelInput, r.Result)
	}

	// send some configs
	configs := []string{
		"interface loopback0",
		"interface loopback0 description tacocat",
		"no interface loopback0",
	}

	_, err = d.SendConfigs(configs)
	if err != nil {
		fmt.Printf("failed to send configs; error: %+v\n", err)
		return
	}
}
