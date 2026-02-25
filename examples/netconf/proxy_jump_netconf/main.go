package main

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	scrapligoconstants "github.com/scrapli/scrapligo/v2/constants"
	scrapligoffi "github.com/scrapli/scrapligo/v2/ffi"
	scrapligonetconf "github.com/scrapli/scrapligo/v2/netconf"
	scrapligooptions "github.com/scrapli/scrapligo/v2/options"
	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
)

const (
	isDarwin = runtime.GOOS == "darwin"

	defaultTimeout = 30 * time.Second

	defaultHostLinux  = "172.20.20.19"
	defaultHostDarwin = "localhost"

	defaultPortLinux  = 22
	defaultPortDarwin = 24022

	netconfPort = 830
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

	n, err := scrapligonetconf.NewNetconf(
		// unlike other examples going by name since we have the config file here
		"srl",
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

	// the simplest thing to do of course is just a "GetConfig"
	result, err := n.GetConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed getting config, error: %v", err))
	}

	// the srl config is *enormous* so just print a little of it...
	fmt.Println(result.Result[0:250])
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
		// now the bits for connecting to the final host, this is always (in this example) to port
		// 830 and to the "real" address (docker address) of the srl box since its from the
		// perspective of the jump host
		scrapligooptions.WithSSH2ProxyJumpHost("172.20.20.16"),
		scrapligooptions.WithSSH2ProxyJumpPort(netconfPort),
		scrapligooptions.WithSSH2ProxyJumpUsername("admin"),
		scrapligooptions.WithSSH2ProxyJumpPassword("NokiaSrl1!"),
	}

	n, err := scrapligonetconf.NewNetconf(
		// note we are still connecting to the bastion on 22 here, then from there to the srl box
		// we are on 830 like youd expect
		defaultHost(),
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

	// the simplest thing to do of course is just a "GetConfig"
	result, err := n.GetConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed getting config, error: %v", err))
	}

	// the srl config is *enormous* so just print a little of it...
	fmt.Println(result.Result[0:250])
}

func main() {
	viaBinTransport()

	fmt.Println()
	fmt.Println()

	viaSSH2Transport()
}
