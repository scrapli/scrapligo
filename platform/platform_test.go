package platform_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/platform"

	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util"
)

var (
	update = flag.Bool( //nolint
		"update",
		false,
		"update the golden files",
	)
	functional = flag.Bool( //nolint
		"functional",
		false,
		"execute functional tests",
	)
	platforms = flag.String( //nolint
		"platforms",
		util.All,
		"comma sep list of platform(s) to target",
	)
	transports = flag.String( //nolint
		"transports",
		util.All,
		"comma sep list of transport(s) to target",
	)
)

func resolveFile(t *testing.T, f string) string {
	f, err := filepath.Abs(fmt.Sprintf("./test-fixtures/%s", f))
	if err != nil {
		t.Fatal(err)
	}

	return f
}

func readFile(t *testing.T, f string) []byte {
	b, err := os.ReadFile(fmt.Sprintf("./test-fixtures/%s", f))
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func writeGolden(t *testing.T, testName string, actualOut []byte) {
	goldenOut := filepath.Join("test-fixtures", "golden", testName+"-out.txt")

	err := os.WriteFile(goldenOut, actualOut, 0o644) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}
}

func prepareDriver(
	t *testing.T,
	testName,
	platformFile,
	payloadFile string,
) (*network.Driver, *transport.File) {
	p, err := platform.NewPlatform(
		resolveFile(t, platformFile),
		"dummy",
		options.WithTransportType(transport.FileTransport),
		options.WithFileTransportFile(resolveFile(t, payloadFile)),
		options.WithTransportReadSize(1),
	)
	if err != nil {
		t.Errorf("%s: encountered error creating Platform, error: %s", testName, err)
	}

	d, err := p.GetNetworkDriver()
	if err != nil {
		t.Errorf(
			"%s: encountered error fetching network driver from platform, error: %s",
			testName,
			err,
		)
	}

	err = d.Channel.Open()
	if err != nil {
		t.Errorf("%s: encountered error opening Channel, error: %s", testName, err)
	}

	fileTransportObj, ok := d.Transport.Impl.(*transport.File)
	if !ok {
		t.Fatal("transport implementation is not Transport File")
	}

	return d, fileTransportObj
}
