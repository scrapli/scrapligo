package main

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
	scrapligoconstants "github.com/scrapli/scrapligo/v2/constants"
	scrapligoffi "github.com/scrapli/scrapligo/v2/ffi"
	scrapligooptions "github.com/scrapli/scrapligo/v2/options"
	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
)

const (
	isDarwin = runtime.GOOS == "darwin"

	defaultTimeout = 30 * time.Second

	defaultPlatform = scrapligocli.NokiaSrlinux

	defaultHostLinux  = "172.20.20.19"
	defaultHostDarwin = "localhost"

	defaultPortLinux  = 22
	defaultPortDarwin = 24022
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

func viaBinTransport() {
	defer func() {
		// will do nothing/no checking unless `LIBSCRAPLI_DEBUG=1` set or using a debug (non
		// release) build of libscrapli
		err := scrapligoffi.AssertNoLeaks()
		if err != nil {
			panic(err)
		}
	}()

	// for proxy-jumping with the bin transport (meaning /bin/ssh literally) we just do
	// normal proxyjump stuff in a config file, then ensure we pass that file. we pick
	// based on linux/darwin here since in darwin we'll hit the exposed ports vs being able to
	// go straight to the ips on the bridge on a linux box
	_, thisFileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed getting path to this example file")
	}

	dir := filepath.Dir(thisFileName)

	var sshConfigFile string

	if runtime.GOOS == scrapligoconstants.Darwin {
		sshConfigFile = filepath.Join(dir, "ssh_config_darwin")
	} else {
		sshConfigFile = filepath.Join(dir, "ssh_config_linux")
	}

	// being lazy and just using one big context for the whole example
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	opts := []scrapligooptions.Option{
		scrapligooptions.WithDefinitionFileOrName(
			scrapligoutil.GetEnvStrOrDefault(
				"SCRAPLI_PLATFORM",
				defaultPlatform.String(),
			),
		),
		// we dont pass port here like other examples since its via the jumper host and we dont
		// need to faff w/ the docker mac nat stuff
		scrapligooptions.WithUsername(
			scrapligoutil.GetEnvStrOrDefault("SCRAPLI_USERNAME", "admin"),
		),
		scrapligooptions.WithPassword(
			scrapligoutil.GetEnvStrOrDefault("SCRAPLI_PASSWORD", "NokiaSrl1!"),
		),
		scrapligooptions.WithBinTransportSSHConfigFile(sshConfigFile),
	}

	c, err := scrapligocli.NewCli(
		// unlike other examples going by name since we have the config file here
		"srl",
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

	result, err := c.SendInput(ctx, "show version")
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())
}

func viaSSH2Transport() {
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

	// libssh2 is a little different -- we setup the connection to the *bastion host* like how we
	// normally setup the connection, then under the transport options you can specify how to
	// connect to the final host
	opts := []scrapligooptions.Option{
		scrapligooptions.WithDefinitionFileOrName(
			scrapligoutil.GetEnvStrOrDefault(
				"SCRAPLI_PLATFORM",
				defaultPlatform.String(),
			),
		),
		// this is to the *bastion host* in the libssh2 case (jumper)
		scrapligooptions.WithUsername(
			scrapligoutil.GetEnvStrOrDefault("SCRAPLI_USERNAME", "scrapli-pw"),
		),
		scrapligooptions.WithPassword(
			scrapligoutil.GetEnvStrOrDefault("SCRAPLI_PASSWORD", "scrapli-123-pw"),
		),
		scrapligooptions.WithTransportSSH2(),
		scrapligooptions.WithPort(
			uint16(scrapligoutil.GetEnvIntOrDefault("SCRAPLI_PORT", defaultPort())), //nolint: gosec
		),
		// now the bits for connecting to the final host
		scrapligooptions.WithSSH2ProxyJumpHost("172.20.20.16"),
		scrapligooptions.WithSSH2ProxyJumpUsername("admin"),
		scrapligooptions.WithSSH2ProxyJumpPassword("NokiaSrl1!"),
	}

	c, err := scrapligocli.NewCli(
		defaultHost(),
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

	result, err := c.SendInput(ctx, "show version")
	if err != nil {
		panic(fmt.Sprintf("failed sending input, error: %v", err))
	}

	fmt.Println(result.Result())
}

func main() {
	viaBinTransport()

	fmt.Println()
	fmt.Println()

	viaSSH2Transport()
}
