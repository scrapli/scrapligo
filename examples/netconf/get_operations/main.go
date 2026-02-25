package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	scrapligoffi "github.com/scrapli/scrapligo/v2/ffi"
	scrapligonetconf "github.com/scrapli/scrapligo/v2/netconf"
	scrapligooptions "github.com/scrapli/scrapligo/v2/options"
	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
)

const (
	isDarwin = runtime.GOOS == "darwin"

	defaultTimeout = 30 * time.Second

	defaultHostLinux          = "172.20.20.16"
	defaultHostDarwin         = "localhost"
	defaultNetopeerHostLinux  = "172.20.20.18"
	defaultNetopeerHostDarwin = "localhost"

	defaultPortLinux          = 830
	defaultPortDarwin         = 21830
	defaultNetopeerPortLinux  = 830
	defaultNetopeerPortDarwin = 23830
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
			scrapligoutil.GetEnvStrOrDefault("SCRAPLI_USERNAME", "admin"),
		),
		scrapligooptions.WithPassword(
			scrapligoutil.GetEnvStrOrDefault("SCRAPLI_PASSWORD", "NokiaSrl1!"),
		),
	}

	return host, opts
}

func defaultNetopeerHost() string {
	if isDarwin {
		return defaultNetopeerHostDarwin
	}

	return defaultNetopeerHostLinux
}

func defaultNetopeerPort() int {
	if isDarwin {
		return defaultNetopeerPortDarwin
	}

	return defaultNetopeerPortLinux
}

func getNetopeerOptions() (string, []scrapligooptions.Option) {
	host := defaultNetopeerHost()

	opts := []scrapligooptions.Option{
		scrapligooptions.WithPort(
			uint16( //nolint:gosec
				scrapligoutil.GetEnvIntOrDefault("SCRAPLI_PORT", defaultNetopeerPort()),
			),
		),
		scrapligooptions.WithUsername("root"),
		scrapligooptions.WithPassword("password"),
	}

	return host, opts
}

func main() {
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

	// the simplest thing to do of course is just a "GetConfig"
	result, err := n.GetConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed getting config, error: %v", err))
	}

	// the srl config is *enormous* so just print a little of it...
	fmt.Println(result.Result[0:250])

	// we can also do "get" rpcs of course... here we'll just provide some simple filter for
	// snagging acl info; you can provide a filter to get_config in the same fashion.
	// default filter type is subtree, srlinux doesnt support xpath, so cant check that here,
	// but you can go head over to the functional tests to see that... but... basically just set
	// the filter_type and then pass a valid xpath filter
	result, err = n.Get(
		ctx,
		scrapligonetconf.WithFilter(`<acl xmlns="urn:nokia.com:srlinux:acl:acl"></acl>`),
	)
	if err != nil {
		panic(fmt.Sprintf("failed running get, error: %v", err))
	}

	fmt.Println(result.Result[0:250])

	// you may also just wanna snag the schema:
	result, err = n.GetSchema(ctx, "ietf-yang-types")
	if err != nil {
		panic(fmt.Sprintf("failed running get schema, error: %v", err))
	}

	// also large output, you get the idea...
	fmt.Println(result.Result[0:250])

	// for get-data we can use the netopeer server instead
	host, opts = getNetopeerOptions()

	n2, err := scrapligonetconf.NewNetconf(
		host,
		opts...,
	)
	if err != nil {
		panic(fmt.Sprintf("failed creating netconf object, error: %v", err))
	}

	_, err = n2.Open(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed opening netconf object, error: %v", err))
	}

	defer func() {
		// once opened always make sure to defer closing! if you dont you will leak memory :)
		_, _ = n2.Close(ctx)
	}()

	// GetData works as you'd expect, basically same as other gets...
	result, err = n2.GetData(
		ctx,
		scrapligonetconf.WithFilter(`<system xmlns="urn:some:data"></system>`),
	)
	if err != nil {
		panic(fmt.Sprintf("failed running get data, error: %v", err))
	}

	fmt.Println(result.Result)
}
