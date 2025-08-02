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

	ctrlC = "\x03"
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

func loggedInCallback(_ context.Context, _ *scrapligocli.Cli) error {
	fmt.Println("a user has logged in!")

	return nil
}

func loggedOutCallback(ctx context.Context, c *scrapligocli.Cli) error {
	fmt.Println("a user has logged out!")

	_, err := c.SendInput(
		ctx,
		ctrlC,
		scrapligocli.WithInputHandling(scrapligocli.InputHandlingIgnore),
		scrapligocli.WithRequestedMode("bash"),
	)

	return err
}

func runReadWithCallbacks(ctx context.Context, cancel context.CancelFunc) {
	defer func() { //nolint:contextcheck
		// will do nothing/no checking unless `LIBSCRAPLI_DEBUG=1` set or using a debug (non
		// release) build of libscrapli
		err := scrapligoffi.AssertNoLeaks()
		if err != nil {
			panic(err)
		}
	}()

	defer cancel()

	host, opts := getOptions()

	c, err := scrapligocli.NewCli( //nolint: contextcheck
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

	// ReadWithCallbacks doesnt support "normal" options, so we just enter mode before hand since
	// we need to be in bash
	_, _ = c.EnterMode(ctx, "bash")

	// note: many real world use cases of this you would want to include a timeout of 0 here
	// so that you *never* actually timeout and instead rely on your callbacks to know when to
	// stop reading from the device
	result, err := c.ReadWithCallbacks(
		ctx,
		"tail -f /var/log/messages",
		scrapligocli.NewReadCallback(
			"user-logged-in",
			loggedInCallback,
			scrapligocli.WithContains("Starting session"),
			scrapligocli.WithOnce(),
		),
		scrapligocli.NewReadCallback(
			"user-logged-out",
			loggedOutCallback,
			scrapligocli.WithContains("disconnected by user"),
			scrapligocli.WithOnce(),
			scrapligocli.WithCompletes(),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed running read with callbacks, error: %v", err))
	}

	fmt.Println(result.Result())
}

func runTrigger(ctx context.Context) {
	defer func() { //nolint:contextcheck
		// will do nothing/no checking unless `LIBSCRAPLI_DEBUG=1` set or using a debug (non
		// release) build of libscrapli
		err := scrapligoffi.AssertNoLeaks()
		if err != nil {
			panic(err)
		}
	}()

	host, opts := getOptions()

	c, err := scrapligocli.NewCli( //nolint: contextcheck
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

	// we really just wanna log in/out but we may as well do something while we are in there
	result, err := c.SendInput(
		ctx,
		"show version",
	)
	if err != nil {
		panic(fmt.Sprintf("failed running read with callbacks, error: %v", err))
	}

	fmt.Println(result.Result())
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)

	go runReadWithCallbacks(ctx, cancel)

	go runTrigger(ctx)

	// wait till the read with callbacks goroutine tells us its done
	<-ctx.Done()
}
