package main

import (
	"fmt"
	"log"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/options"
)

func myLoggerFunc(x ...interface{}) {
	// not sure why you would want to do this, but this is just an example logger function you
	// could use :)
	if len(x) > 0 {
		// pro tip... its almost certainly going to be a string, but this is just to show how
		// you can add logger functions so who cares!
		fmt.Printf("got a log item of type %T\n", x[0])
	}
}

func main() {
	li, err := logging.NewInstance(
		logging.WithLevel(logging.Debug),
		logging.WithLogger(log.Print),
		logging.WithLogger(myLoggerFunc),
	)
	if err != nil {
		fmt.Printf("failed to logging instance; error: %+v\n", err)

		return
	}

	d, err := generic.NewDriver(
		"sandbox-iosxe-latest-1.cisco.com",
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername("developer"),
		options.WithAuthPassword("C1sco12345"),
		// the options.WithLogger applies the logging instance created above to our driver
		options.WithLogger(li),
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

	fmt.Printf("got some output: %s\n\n\n", r.Result)
}
