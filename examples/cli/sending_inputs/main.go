package main

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
	scrapligoffi "github.com/scrapli/scrapligo/v2/ffi"
	scrapligooptions "github.com/scrapli/scrapligo/v2/options"
	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
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
		scrapligooptions.WithDefinitionFileOrName(
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

	_, thisFileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed getting path to this example file")
	}

	dir := filepath.Dir(thisFileName)

	inputsFromFile := filepath.Join(dir, "inputs_to_send")

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

	// send a single "input" at the "default mode" (normally privileged exec or similar)
	result, err := c.SendInput(ctx, "show version")
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	// the result object returned holds info about the operation -- start/end/duration, the
	// input(s) sent, the result, the raw result (as in before ascii/ansii cleaning), and a few
	// other things. it has a reasonable __str__ method, so printing it should give you some
	// something to look at
	fmt.Println(result)

	// but if you want to just see the result itself you can do like so:
	fmt.Println(result.Result())

	// theres a plural method for... sending multiple inputs, shock!
	results, err := c.SendInputs(ctx, []string{"show version", "show version"})
	if err != nil {
		panic(fmt.Sprintf("failed sending inputs, error: %v", err))
	}

	// result will print a joined result
	fmt.Println(results.Result())

	// there is also a from_file method to send inputs from a file if you want
	results, err = c.SendInputsFromFile(ctx, inputsFromFile)
	if err != nil {
		panic(fmt.Sprintf("failed sending inputs, error: %v", err))
	}

	fmt.Println(results.Result())
}
