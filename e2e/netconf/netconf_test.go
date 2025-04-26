package netconf_test

import (
	"context"
	"os"
	"runtime"
	"strings"
	"testing"

	scrapligocli "github.com/scrapli/scrapligo/cli"
	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestMain(m *testing.M) {
	scrapligotesthelper.Flags()

	os.Exit(m.Run())
}

func getTransports() []string {
	return []string{
		"bin",
		"ssh2",
	}
}

func shouldSkipPlatform(platform string) bool {
	if *scrapligotesthelper.Platforms == "all" {
		return false
	}

	platforms := strings.Split(*scrapligotesthelper.Platforms, ",")

	for _, platformName := range platforms {
		if platformName == platform {
			return false
		}
	}

	return true
}

func shouldSkipTransport(transport string) bool {
	if *scrapligotesthelper.Transports == "all" {
		return false
	}

	transports := strings.Split(*scrapligotesthelper.Transports, ",")

	for _, transportName := range transports {
		if transportName == transport {
			return false
		}
	}

	return true
}

func shouldSkip(platform, transport string) bool {
	if shouldSkipPlatform(platform) {
		return true
	}

	if shouldSkipTransport(transport) {
		return true
	}

	return false
}

func getNetconf(t *testing.T, platform, transportName string) *scrapligonetconf.Netconf {
	t.Helper()

	var host string

	var opts []scrapligooptions.Option

	switch transportName {
	case "bin":
		opts = append(
			opts,
			scrapligooptions.WithTransportBin(),
		)
	case "ssh2":
		opts = append(
			opts,
			scrapligooptions.WithTransportSSH2(),
		)
	default:
		t.Fatal("unsupported transport name")
	}

	if platform == scrapligocli.NokiaSrl.String() {
		opts = append(
			opts,
			scrapligooptions.WithUsername("admin"),
			scrapligooptions.WithPassword("NokiaSrl1!"),
		)

		if runtime.GOOS == "darwin" {
			host = "localhost"

			opts = append(
				opts,
				scrapligooptions.WithPort(21830),
			)
		} else {
			host = "172.20.20.16"
		}
	} else {
		opts = append(
			opts,
			scrapligooptions.WithUsername("netconf-admin"),
			scrapligooptions.WithPassword("admin"),
		)

		if runtime.GOOS == "darwin" {
			host = "localhost"

			opts = append(
				opts,
				scrapligooptions.WithPort(22830),
			)
		} else {
			host = "172.20.20.17"
		}
	}

	n, err := scrapligonetconf.NewNetconf(
		host,
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return n
}

func closeNetconf(t *testing.T, n *scrapligonetconf.Netconf) {
	t.Helper()

	_, err := n.Close(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
