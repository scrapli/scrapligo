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
	sessionRecorderPath := filepath.Join(dir, "session_record.log")

	// being lazy and just using one big context for the whole example
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	host, opts := getOptions()

	opts = append(
		opts,
		scrapligooptions.WithSessionRecorderPath(sessionRecorderPath),
	)

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

	// we'll just send stuff to have something to look at
	result, err := c.SendInput(ctx, "show version")
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())
}
