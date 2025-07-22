package main

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligoutil "github.com/scrapli/scrapligo/util"
)

const (
	isDarwin = runtime.GOOS == "darwin"

	defaultTimeout = 30 * time.Second

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

	return scrapligoutil.GetEnvStrOrDefault("SCRAPLI_HOST", defaultHost()), opts
}

func main() {
	// being lazy and just using one big context for the whole example
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	_, thisFileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed getting path to this example file")
	}

	dir := filepath.Dir(thisFileName)
	definitionPath := filepath.Join(dir, "foo_bar.yaml")

	host, opts := getOptions()

	c, err := scrapligocli.NewCli(
		// this is exactly the same as the upstream definition but just doing this to show that
		// you can load up any yaml definition and dont necessarily need to rely on the upstream
		// stuff in scrapli_definitions
		definitionPath,
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
		_, _ = c.Close(ctx)
	}()

	result, err := c.SendInput(ctx, "show version")
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())
}
