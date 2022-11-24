package generic_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/util"

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

func writeGolden(t *testing.T, testName string, actualIn []byte, actualOut string) {
	goldenOut := filepath.Join("test-fixtures", "golden", testName+"-out.txt")
	goldenIn := filepath.Join("test-fixtures", "golden", testName+"-in.txt")

	err := os.WriteFile(goldenOut, []byte(actualOut), 0o644) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(goldenIn, actualIn, 0o644) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}
}

func prepareDriver(
	t *testing.T,
	testName,
	payloadFile string,
) (*generic.Driver, *transport.File) {
	d, err := generic.NewDriver(
		"dummy",
		options.WithTransportType(transport.FileTransport),
		options.WithFileTransportFile(resolveFile(t, payloadFile)),
		options.WithTransportReadSize(1),
		options.WithReadDelay(0),
		options.WithFailedWhenContains([]string{
			"% Ambiguous command",
			"% Incomplete command",
			"% Invalid input detected",
			"% Unknown command",
		}),
	)
	if err != nil {
		t.Errorf("%s: encountered error creating generic Driver, error: %s", testName, err)
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
