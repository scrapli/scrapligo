package cli_test

import (
	"os"
	"runtime"
	"strings"
	"testing"

	scrapligocli "github.com/scrapli/scrapligo/cli"
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
		"telnet",
	}
}

func shouldSkipPlatform(platform string) bool {
	if *scrapligotesthelper.Platforms == "all" {
		return false
	}

	platforms := strings.Split(*scrapligotesthelper.Platforms, ",")

	for _, platformName := range platforms {
		if platformName == platform {
			return true
		}
	}

	return false
}

func shouldSkipTransport(transport string) bool {
	if *scrapligotesthelper.Transports == "all" {
		return false
	}

	transports := strings.Split(*scrapligotesthelper.Transports, ",")

	for _, transportName := range transports {
		if transportName == transport {
			return true
		}
	}

	return false
}

func shouldSkip(platform, transport string) bool {
	if shouldSkipPlatform(platform) {
		return true
	}

	if shouldSkipTransport(transport) {
		return true
	}

	if transport == "telnet" && platform != scrapligocli.AristaEos.String() {
		// now we just check against telnet, since we only run that against eos for now
		return true
	}

	return false
}

func getDriver(t *testing.T, platform, transportName string) *scrapligocli.Driver {
	var host string

	opts := []scrapligooptions.Option{
		scrapligooptions.WithUsername("admin"),
	}

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
	case "telnet":
		opts = append(
			opts,
			scrapligooptions.WithTransportTelnet(),
		)
	default:
		t.Fatal("unsupported transport name")
	}

	if platform == scrapligocli.NokiaSrl.String() {
		opts = append(
			opts,
			scrapligooptions.WithPassword("admin"),
			scrapligooptions.WithPassword("NokiaSrl1!"),
		)

		if runtime.GOOS == "darwin" {
			host = "localhost"

			opts = append(
				opts,
				scrapligooptions.WithPort(21022),
			)
		} else {
			host = "172.20.20.16"
		}
	} else {
		opts = append(
			opts,
			scrapligooptions.WithPassword("admin"),
			scrapligooptions.WithLookupKeyValue("enable", "libscrapli"),
		)

		var port uint16

		if runtime.GOOS == "darwin" {
			host = "localhost"
			port = 22022
		} else {
			host = "172.20.20.17"
		}

		if transportName == "telnet" {
			port++
		}

		opts = append(
			opts,
			scrapligooptions.WithPort(port),
		)
	}

	d, err := scrapligocli.NewDriver(
		platform,
		host,
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func closeDriver(t *testing.T, d *scrapligocli.Driver) {
	_ = t

	d.Close()
}
