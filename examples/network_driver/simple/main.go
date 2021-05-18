package main

import (
	"bytes"
	"fmt"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/driver/core"
)

const commandsFile = "commandsfile"

func main() {
	var channelLog bytes.Buffer

	d, err := core.NewCoreDriver(
		"localhost",
		"cisco_iosxe",
		base.WithPort(21022),
		base.WithAuthStrictKey(false),
		base.WithAuthUsername("vrnetlab"),
		base.WithAuthPassword("VR-netlab9"),
		base.WithAuthSecondary("VR-netlab9"),
		base.WithChannelLog(&channelLog),
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

	// fetch the prompt
	prompt, err := d.GetPrompt()
	if err != nil {
		fmt.Printf("failed to get prompt; error: %+v\n", err)
	} else {
		fmt.Printf("found prompt: %s\n", prompt)
	}

	// send some commands from a file
	mr, err := d.SendCommandsFromFile(commandsFile)
	if err != nil {
		fmt.Printf("failed to send commands from file; error: %+v\n", err)
		return
	}
	for _, r := range mr.Responses {
		fmt.Printf("sent command '%s', output received:\n %s\n", r.ChannelInput, r.Result)
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

	err = d.Close()
	if err != nil {
		fmt.Printf("failed to close driver; error: %+v\n", err)
	}

	b := make([]byte, 65535)
	_, _ = channelLog.Read(b)
	fmt.Printf("READ: %s", b)
}
