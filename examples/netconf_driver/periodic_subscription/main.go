package main

import (
	"fmt"
	"time"

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

	xpath := "/interfaces/interface[name=\"GigabitEthernet1\"]/statistics"
	period := 1000
	resp, err := d.EstablishPeriodicSubscription(xpath, period)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Subscription ID: ", resp.SubscriptionID)

	for {
		fmt.Println("Checking for messages...")
		messages := d.GetSubscriptionMessages(resp.SubscriptionID)
		for _, event := range messages {
			fmt.Println("Event Received:")
			fmt.Println(string([]byte(event)))
		}
		time.Sleep(1 * time.Second)
	}
}
