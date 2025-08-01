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

	defaultHostLinux  = "172.20.20.16"
	defaultHostDarwin = "localhost"

	defaultPortLinux  = 830
	defaultPortDarwin = 21830
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

	// we can lock the config before doing things if we want
	result, err := n.Lock(
		ctx,
		scrapligonetconf.WithDatastore(scrapligonetconf.DatastoreTypeCandidate),
	)
	if err != nil {
		panic(fmt.Sprintf("failed locking datastore, error: %v", err))
	}

	fmt.Println(result.Result)

	// and push a valid config of course
	result, err = n.EditConfig(
		ctx,
		`
  		<system xmlns="urn:nokia.com:srlinux:general:system">
            <name xmlns="urn:nokia.com:srlinux:chassis:system-name">
                <host-name>foozzzBaaaaAR</host-name>
            </name>
        </system>
		`,
		scrapligonetconf.WithDatastore(scrapligonetconf.DatastoreTypeCandidate),
	)
	if err != nil {
		panic(fmt.Sprintf("failed editing config, error: %v", err))
	}

	fmt.Println(result.Result)

	result, err = n.Commit(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed committing config, error: %v", err))
	}

	fmt.Println(result.Result)

	// annnnd unlock
	result, err = n.Unlock(
		ctx,
		scrapligonetconf.WithDatastore(scrapligonetconf.DatastoreTypeCandidate),
	)
	if err != nil {
		panic(fmt.Sprintf("failed unlocking datastore, error: %v", err))
	}

	fmt.Println(result.Result)

	// well put it back just in case using this w/ testing so we didnt make any change
	_, err = n.EditConfig(
		ctx,
		`
  		<system xmlns="urn:nokia.com:srlinux:general:system">
            <name xmlns="urn:nokia.com:srlinux:chassis:system-name">
                <host-name>srl</host-name>
            </name>
        </system>
		`,
		scrapligonetconf.WithDatastore(scrapligonetconf.DatastoreTypeCandidate),
	)
	if err != nil {
		panic(fmt.Sprintf("failed editing config, error: %v", err))
	}

	_, err = n.Commit(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed committing config, error: %v", err))
	}
}
