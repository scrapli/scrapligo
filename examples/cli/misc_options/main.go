package main

import (
	"context"
	"errors"
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

	// WithRetainTrailingPrompt will... retain the prompt after the output from the command you
	// send. this is False by default since normally you'd want just the commands output but
	// obviously sometimes you may want to see the prompt too
	result, err := c.SendInput(ctx, "show version", scrapligocli.WithRetainTrailingPrompt())
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())

	// just for clarity
	fmt.Println()
	fmt.Println()

	// we can also retain the input
	result, err = c.SendInput(ctx, "show version", scrapligocli.WithRetainInput())
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())

	// just for clarity
	fmt.Println()
	fmt.Println()

	// and of course these can be combined
	result, err = c.SendInput(
		ctx,
		"show version",
		scrapligocli.WithRetainTrailingPrompt(),
		scrapligocli.WithRetainInput(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())

	// just for clarity
	fmt.Println()
	fmt.Println()

	// and we can time stuff out using the context as you'd expect
	fastCtx, fastCtxCancel := context.WithTimeout(ctx, 30*time.Nanosecond) //nolint: mnd
	defer fastCtxCancel()

	_, err = c.SendInput(fastCtx, "show version")
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			panic(fmt.Sprintf("should have been a deadline exceeded error but got: %v", err))
		}
	}

	// the last option is "input handling" -- there are a few options here:
	// - Exact
	// - Fuzzy
	// - Ignore
	// the gist here is that scrapli "looks" for your inputs before sending the return --
	// historically scrapli has looked for the *exact* input. this is *usually* good, but there
	// are some places where the input you send is not what is reflected in the session; things
	// like banners or "vi-like" input modes, or when a device writes \x08 (backspaces) when your
	// input is going lonver than the terminal width. so, nowadays scrapli "fuzzily" matches the
	// input -- meaning that as long as all the characters you send are in the output in the same
	// order (but if you send "foo" then "f X o X o" would be allowed). lastly, you can *ignore*
	// the input... but you shouldnt do this. this is really just used in netconf operations.
	_, _ = c.SendInput(
		ctx,
		"show version",
		scrapligocli.WithInputHandling(scrapligocli.InputHandlingExact),
	)
}
