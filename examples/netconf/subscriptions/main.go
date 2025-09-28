package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	scrapligoffi "github.com/scrapli/scrapligo/ffi"
	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligoutil "github.com/scrapli/scrapligo/util"
)

const (
	isDarwin = runtime.GOOS == "darwin"

	defaultTimeout = 30 * time.Second

	defaultHostLinux  = "172.20.20.18"
	defaultHostDarwin = "localhost"

	defaultPortLinux  = 830
	defaultPortDarwin = 23830

	tryIosxe      = false
	iosxeHost     = ""
	iosxePort     = 830
	iosxeUsername = ""
	iosxePassword = ""
)

func defaultHost() string {
	if isDarwin {
		return defaultHostDarwin
	}

	return defaultHostLinux
}

func defaultPort() int {
	if isDarwin {
		return defaultPortDarwin
	}

	return defaultPortLinux
}

func getOptions() (string, []scrapligooptions.Option) {
	host := scrapligoutil.GetEnvStrOrDefault(
		"SCRAPLI_HOST",
		defaultHost(),
	)

	opts := []scrapligooptions.Option{
		scrapligooptions.WithPort(
			uint16(scrapligoutil.GetEnvIntOrDefault("SCRAPLI_PORT", defaultPort())), //nolint:gosec
		),
		scrapligooptions.WithUsername(
			scrapligoutil.GetEnvStrOrDefault("SCRAPLI_USERNAME", "root"),
		),
		scrapligooptions.WithPassword(
			scrapligoutil.GetEnvStrOrDefault("SCRAPLI_PASSWORD", "password"),
		),
	}

	return host, opts
}

func notifications() {
	defer func() {
		// will do nothing/no checking unless `LIBSCRAPLI_DEBUG=1` set or using a debug (non
		// release) build of libscrapli
		err := scrapligoffi.AssertNoLeaks()
		if err != nil {
			panic(err)
		}
	}()

	// being lazy and just using one big context for the whole example
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	host, opts := getOptions()

	n, err := scrapligonetconf.NewNetconf(
		host,
		opts...,
	)
	if err != nil {
		panic(fmt.Sprintf("failed creating netconf object, error: %v", err))
	}

	_, err = n.Open(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed opening netconf object, error: %v", err))
	}

	defer func() {
		// once opened always make sure to defer closing! if you dont you will leak memory :)
		_, _ = n.Close(ctx)
	}()

	// because there are a zillion variants to how to setup subscriptions -- i.e. create vs
	// establish then differing ways to setup the payload based on the rfc that is followed
	// scrapli decided... nope. you can just send what you need to create your subscription
	// however makes sense for your server. here we'll just do a very simple example.
	result, err := n.RawRPC(
		ctx,
		`
		<create-subscription xmlns="urn:ietf:params:xml:ns:netconf:notification:1.0">
		</create-subscription>
		`,
	)
	if err != nil {
		panic(fmt.Sprintf("failed sending rpc, error: %v", err))
	}

	fmt.Println(result.Result)

	_, err = n.GetNextNotification()
	if err == nil {
		// expected an error as the netopeer server churns out a notification every 3s or
		// something for the "boring counter" object
		panic("we expected no notifications yet")
	}

	// enough to ensure we can fetch a notification
	time.Sleep(10 * time.Second) //nolint: mnd

	notif, err := n.GetNextNotification()
	if err != nil {
		panic("womp womp, no notifications to snag sadly...t")
	}

	fmt.Println(notif)
}

func subscriptions() {
	defer func() {
		// will do nothing/no checking unless `LIBSCRAPLI_DEBUG=1` set or using a debug (non
		// release) build of libscrapli
		err := scrapligoffi.AssertNoLeaks()
		if err != nil {
			panic(err)
		}
	}()

	// being lazy and just using one big context for the whole example
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// the following assumes some standard iosxe device (with netconf running ofc) like iosxe
	// cat8k always on sandbox. it shows using *establish-subscription* and then fetching
	// subscription messages by the subscription id returned in the establish rpc.
	n, err := scrapligonetconf.NewNetconf(
		iosxeHost,
		scrapligooptions.WithPort(iosxePort),
		scrapligooptions.WithUsername(iosxeUsername),
		scrapligooptions.WithPassword(iosxePassword),
	)
	if err != nil {
		panic(fmt.Sprintf("failed creating netconf object, error: %v", err))
	}

	_, err = n.Open(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed opening netconf object, error: %v", err))
	}

	defer func() {
		// once opened always make sure to defer closing! if you dont you will leak memory :)
		_, _ = n.Close(ctx)
	}()

	result, err := n.RawRPC(
		ctx,
		//nolint: lll
		`
		<establish-subscription xmlns="urn:ietf:params:xml:ns:yang:ietf-event-notifications" xmlns:yp="urn:ietf:params:xml:ns:yang:ietf-yang-push">
            <stream>yp:yang-push</stream>
            <yp:xpath-filter>/mdt-oper:mdt-oper-data/mdt-subscriptions</yp:xpath-filter>
            <yp:period>1000</yp:period>
        </establish-subscription>
		`,
	)
	if err != nil {
		panic(fmt.Sprintf("failed sending rpc, error: %v", err))
	}

	fmt.Println(result.Result)

	subscriptionID, err := n.GetSubscriptionID(result.Result)
	if err != nil {
		panic(fmt.Sprintf("failed getting subscription id, error: %v", err))
	}

	fmt.Println("subscription id > ", subscriptionID)

	for {
		time.Sleep(3 * time.Second) //nolint: mnd

		notif, err := n.GetNextSubscription(subscriptionID)
		if err == nil {
			fmt.Println(notif)

			return
		}

		fmt.Println("sad panda, no subscription messages... *yet*")
	}
}

func main() {
	notifications()

	if tryIosxe == false {
		return
	}

	subscriptions()
}
