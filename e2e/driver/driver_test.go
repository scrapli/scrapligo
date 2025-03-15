package driver_test

import (
	"os"
	"runtime"
	"strings"
	"testing"

	scrapligodriver "github.com/scrapli/scrapligo/driver"
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

func shouldSkip(platform, transport string) bool {
	if scrapligotesthelper.Platforms != nil {
		var skipPlatform bool

		platforms := strings.Split(*scrapligotesthelper.Platforms, ",")

		for _, platformName := range platforms {
			if platformName == platform {
				skipPlatform = true
			}
		}

		if skipPlatform {
			return false
		}
	}

	if scrapligotesthelper.Transports != nil {
		var skipTransport bool

		transports := strings.Split(*scrapligotesthelper.Transports, ",")

		for _, transportName := range transports {
			if transportName == transport {
				skipTransport = true
			}
		}

		if skipTransport {
			return false
		}
	}

	if transport == "telnet" && platform == scrapligodriver.NokiaSrl.String() {
		// no telnet on srl node
		return false
	}

	return true
}

func getDriver(t *testing.T, platform, transportName string) *scrapligodriver.Driver {
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

	if platform == string(scrapligodriver.NokiaSrl) {
		opts = append(
			opts,
			scrapligooptions.WithPassword("admin"),
		)

		if runtime.GOOS == "darwin" {
			host = "localhost"

			opts = append(
				opts,
				scrapligooptions.WithPassword("NokiaSrl1!"),
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
			port = 22
		}

		if transportName == "telnet" {
			// TODO i think that we have to add the annoying "send return when you dont see shit"
			//  to zig in session auth for telnet stuff
			port++
		}

		opts = append(
			opts,
			scrapligooptions.WithPort(port),
		)
	}

	d, err := scrapligodriver.NewDriver(
		platform,
		host,
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func closeDriver(t *testing.T, d *scrapligodriver.Driver) {
	_ = t

	d.Close()
}
