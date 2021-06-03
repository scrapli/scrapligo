package main

import (
	"bytes"
	"fmt"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/driver/core"
)

func main() {
	// WithChannelLog accepts an io.Writer type object, create and pass it to the driver creation
	var channelLog bytes.Buffer

	// use the NewCoreDriver factory and pass in a platform argument
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

	prompt, err := d.GetPrompt()
	if err != nil {
		fmt.Printf("failed to get prompt; error: %+v\n", err)
		return
	}
	fmt.Printf("found prompt: %s\n", prompt)

	err = d.Close()
	if err != nil {
		fmt.Printf("failed to close driver; error: %+v\n", err)
	}

	// We can then read and print out the channel log data like normal
	b := make([]byte, 0, 65535)
	_, _ = channelLog.Read(b)
	fmt.Printf("Channel log output:\n%s", b)
}
