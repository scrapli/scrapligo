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

func getOptions() (string, string, []scrapligooptions.Option) { //nolint: gocritic
	platform := scrapligoutil.GetEnvStrOrDefault(
		"SCRAPLI_PLATFORM",
		defaultPlatform.String(),
	)

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

	return platform, host, opts
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

	platform, host, opts := getOptions()

	c, err := scrapligocli.NewCli(
		platform,
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

	// you can manually request to "enter" a mode -- in this case we can go into configuration
	// mode. note that if you were to issue a subsequent send_input *without specifying*
	// configuration as the requested mode you would be "dropped" back into exec (the default
	// preferred mode)!	result, err := c.SendPromptedInput(
	_, err = c.EnterMode(ctx, "configuration")
	if err != nil {
		panic(fmt.Sprintf("failed entering mode, error: %v", err))
	}

	result, err := c.GetPrompt(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed getting prompt, error: %v", err))
	}

	// note the "candidate private" -- we are in configuration mode
	fmt.Println(result.Result())

	// just to make the outputs clearer
	fmt.Println()
	fmt.Println()

	// to illustrate that we auto try to send inputs from the default mode we can issue a show
	// version, retaining the trailing prompt to verify we are in fact no longer in config mode
	result, err = c.SendInput(ctx, "show version", scrapligocli.WithRetainTrailingPrompt())
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())
}
