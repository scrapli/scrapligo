package channel_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/logging"
	"github.com/scrapli/scrapligo/transport"
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

func writeGolden(t *testing.T, testName string, actualIn, actualOut []byte) {
	goldenOut := filepath.Join("test-fixtures", "golden", testName+"-out.txt")
	goldenIn := filepath.Join("test-fixtures", "golden", testName+"-in.txt")

	err := os.WriteFile(goldenOut, actualOut, 0o644) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(goldenIn, actualIn, 0o644) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}
}

func prepareChannel(
	t *testing.T,
	testName,
	payloadFile string,
) (*channel.Channel, *transport.File) {
	l, _ := logging.NewInstance()

	transportObj, err := transport.NewTransport(
		l,
		"dummy",
		transport.FileTransport,
		options.WithFileTransportFile(resolveFile(t, payloadFile)),
		options.WithTransportReadSize(1),
	)
	if err != nil {
		t.Errorf("%s: encountered error creating File Transport, error: %s", testName, err)
	}

	c, err := channel.NewChannel(l, transportObj)
	if err != nil {
		t.Errorf("%s: encountered error creating Channel, error: %s", testName, err)
	}

	err = c.Open()
	if err != nil {
		t.Errorf("%s: encountered error opening Channel, error: %s", testName, err)
	}

	fileTransportObj, ok := transportObj.Impl.(*transport.File)
	if !ok {
		t.Fatal("transport implementation is not Transport File")
	}

	return c, fileTransportObj
}
