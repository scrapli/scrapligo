package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
	scrapligoffi "github.com/scrapli/scrapligo/v2/ffi"
	scrapligologging "github.com/scrapli/scrapligo/v2/logging"
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

func stdLogger() {
	// being lazy and just using one big context for the whole example
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	host, opts := getOptions()

	opts = append(
		opts,
		// we can simply pass the "normal" std library logger, but make sure to pass/set a level
		// too otherwise its just warn so you wont see much (or hopefully anything)
		scrapligooptions.WithLogger(log.Default()),
		scrapligooptions.WithLoggerLevel(scrapligologging.Debug),
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

	result, err := c.SendInput(ctx, "show version", scrapligocli.WithRetainTrailingPrompt())
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())
}

func slogLogger() {
	// being lazy and just using one big context for the whole example
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	host, opts := getOptions()

	opts = append(
		opts,
		// or the std library slog logger, but make sure to pass/set a level
		// too otherwise its just warn so you wont see much (or hopefully anything)
		scrapligooptions.WithLogger(
			slog.New(
				slog.NewTextHandler(
					os.Stdout,
					// w/ slog you need to tell scrapli *and* the slog the level you want
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			),
		),
		scrapligooptions.WithLoggerLevel(scrapligologging.Debug),
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

	result, err := c.SendInput(ctx, "show version", scrapligocli.WithRetainTrailingPrompt())
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())
}

func callbackLogger() {
	// being lazy and just using one big context for the whole example
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	host, opts := getOptions()

	opts = append(
		opts,
		// or a custom callback w/ a simple signature
		scrapligooptions.WithLogger(
			func(level scrapligologging.LogLevel, message string) {
				_, _ = fmt.Fprintln(os.Stderr, level, "::", message)
			},
		),
		scrapligooptions.WithLoggerLevel(scrapligologging.Debug),
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

	result, err := c.SendInput(ctx, "show version", scrapligocli.WithRetainTrailingPrompt())
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())
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

	// first w/ the standard/old/normal logger
	stdLogger()

	fmt.Println()
	fmt.Println()

	// annnnd with slog logger
	slogLogger()

	fmt.Println()
	fmt.Println()

	// last option being a custom callback func to do w/e you want
	callbackLogger()
}
