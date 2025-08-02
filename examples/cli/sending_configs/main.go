package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligoutil "github.com/scrapli/scrapligo/util"
)

const (
	isDarwin = runtime.GOOS == "darwin"

	defaultTimeout = 30 * time.Second

	defaultPlatform = scrapligocli.NokiaSrlinux

	defaultHostLinux  = "172.20.20.16"
	defaultHostDarwin = "localhost"

	defaultPortLinux  = 22
	defaultPortDarwin = 21022
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
		scrapligooptions.WithDefintionFileOrName(
			scrapligoutil.GetEnvStrOrDefault(
				"SCRAPLI_PLATFORM",
				defaultPlatform.String(),
			),
		),
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

	c, err := scrapligocli.NewCli(
		host,
		opts...,
	)
	if err != nil {
		panic(fmt.Sprintf("failed creating cli object, error: %v", err))
	}

	_, err = c.Open(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed opening cli object, error: %v", err))
	}

	defer func() {
		// once opened always make sure to defer closing! if you dont you will leak memory :)
		_, _ = c.Close(ctx)
	}()

	// in scrapli (now, not historically) there is really no such thing as a "config" --
	// everything is simply inputs that we send to the device, maybe at different "modes" -- the
	// "mode" possibly being something like "configuration" mode (i.e. the "mode" you get into by
	// doing "config t"). when requesting a different mode, you need to make sure the "mode"
	// exists on the platform definition (see: scrapli_definitions) then entering/exiting modes
	// will be handled for you.
	result, err := c.SendInput(
		ctx,
		"show version",
		scrapligocli.WithRequestedMode("configuration"),
		scrapligocli.WithRetainTrailingPrompt(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	// above we simply sent a "show version", but we did it in the "configuration" mode -- we
	// also added the `WithRetainTrailingPrompt` flag to retain... the trailing prompt -- just so
	// you can see that we are in fact in "config" mode ("enter candidate private" in srlinux).
	fmt.Println(result.Result())

	// if we want to send a "regular" input again, scrapli will automatically do so from the
	// default mode, which is usually exec/privileged_exec. once again, we retain the trailing
	// prompt so you can confirm/see this. note that *how* you "leave" the config mode varies
	// depending on the platform, so check the definition. in this case its "discard now".
	result, err = c.SendInput(
		ctx,
		"show version",
		scrapligocli.WithRetainTrailingPrompt(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())

	// if you want to actually commit/save (obv depending on your device if that is required)
	// you need to actually send the commit/save command yourself
	_, err = c.SendInputs(
		ctx,
		[]string{"set system name host-name foo", "commit now"},
		scrapligocli.WithRequestedMode("configuration"),
	)
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	// and just to confirm for our sanity that it worked...
	result, err = c.SendInput(
		ctx,
		"info system name host-name",
		scrapligocli.WithRetainTrailingPrompt(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())

	// finally, we'll just put that back how we found it...
	_, err = c.SendInputs(
		ctx,
		[]string{"set system name host-name srl", "commit now"},
		scrapligocli.WithRequestedMode("configuration"),
	)
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}
}
