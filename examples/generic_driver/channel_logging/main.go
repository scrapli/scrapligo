package main

import (
	"bytes"
	"fmt"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/options"
)

func main() {
	// WithChannelLog accepts an io.Writer type object, create and pass it to the driver creation
	var channelLog bytes.Buffer

	d, err := generic.NewDriver(
		"sandbox-iosxe-latest-1.cisco.com",
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername("developer"),
		options.WithAuthPassword("C1sco12345"),
		options.WithChannelLog(&channelLog),
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

	// We can then read and print out the channel log data like normal
	b := make([]byte, channelLog.Len())
	_, _ = channelLog.Read(b)
	fmt.Printf("Channel log output:\n%s", b)
}
