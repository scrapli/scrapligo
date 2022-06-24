package netconf_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/transport"

	"github.com/scrapli/scrapligo/platform"
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
) (*netconf.Driver, *transport.File) {
	d, err := netconf.NewDriver(
		"dummy",
		options.WithTransportType(transport.FileTransport),
		options.WithFileTransportFile(resolveFile(t, payloadFile)),
		options.WithTransportReadSize(1),
		options.WithReadDelay(0),
	)
	if err != nil {
		t.Fatalf("%s: encountered error creating network Driver, error: %s", testName, err)
	}

	err = d.Open()
	if err != nil {
		t.Fatalf("%s: encountered error opening netconf Driver, error: %s", testName, err)
	}

	fileTransportObj, ok := d.Transport.Impl.(*transport.File)
	if !ok {
		t.Fatalf("transport implementation is not Transport File")
	}

	return d, fileTransportObj
}

func writeGoldenFunctional(t *testing.T, testName, actualOut string) {
	goldenOut := filepath.Join("test-fixtures", "golden", testName+"-out.txt")

	err := os.WriteFile(goldenOut, []byte(actualOut), 0o644) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}
}

func getNetconfTransportNames() []string {
	ncTransportNames := make([]string, 0)

	for _, transportName := range transport.GetTransportNames() {
		if transportName == transport.TelnetTransport {
			continue
		}

		ncTransportNames = append(ncTransportNames, transportName)
	}

	return ncTransportNames
}

func getNetconfPlatformNames() []string {
	ncPlatformNames := make([]string, 0)

	for _, platformName := range platform.GetPlatformNames() {
		if platformName == platform.NokiaSrl {
			continue
		}

		ncPlatformNames = append(ncPlatformNames, platformName)
	}

	return ncPlatformNames
}

func getFunctionalHostIPPortLinuxOrRemote(
	t *testing.T,
	platformName string,
) (host string, port int) {
	switch platformName {
	case platform.CiscoIosxe:
		host = util.GetEnvStrOrDefault("SCRAPLIGO_CISCO_IOSXE_HOST", "172.20.20.11")

		return host, 830
	case platform.CiscoIosxr:
		host = util.GetEnvStrOrDefault("SCRAPLIGO_CISCO_IOSXR_HOST", "172.20.20.12")

		return host, 830
	case platform.CiscoNxos:
		host = util.GetEnvStrOrDefault("SCRAPLIGO_CISCO_NXOS_HOST", "172.20.20.13")

		return host, 830
	case platform.AristaEos:
		host = util.GetEnvStrOrDefault("SCRAPLIGO_ARISTA_EOS_HOST", "172.20.20.14")

		return host, 830
	case platform.JuniperJunos:
		host = util.GetEnvStrOrDefault("SCRAPLIGO_JUNIPER_JUNOS_HOST", "172.20.20.15")

		return host, 830
	}

	t.Fatalf("failed finding platform host/port info")

	return "", 0
}

func getFunctionalHostIPPort(t *testing.T, platformName string) (host string, port int) {
	osType := runtime.GOOS

	remoteOverride := util.GetEnvIntOrDefault("SCRAPLIGO_NO_HOST_FWD", 0)

	if osType == "linux" || remoteOverride != 0 {
		return getFunctionalHostIPPortLinuxOrRemote(t, platformName)
	}

	// otherwise we are running on darwin w/ local boxen w/ nat setup

	host = "localhost"

	switch platformName {
	case platform.CiscoIosxe:
		return host, 21830
	case platform.CiscoIosxr:
		return host, 22830
	case platform.CiscoNxos:
		return host, 23830
	case platform.AristaEos:
		return host, 24830
	case platform.JuniperJunos:
		return host, 25830
	}

	t.Fatalf("failed finding platform host/port info")

	return "", 0
}

func getFunctionalHostUserPass(t *testing.T, platformName string) (user, pass string) {
	user = util.Admin
	pass = util.Admin

	switch platformName {
	case platform.CiscoIosxe:
		return user, pass
	case platform.CiscoIosxr:
		return "clab", "clab@123"
	case platform.CiscoNxos:
		return user, pass
	case platform.AristaEos:
		return user, pass
	case platform.JuniperJunos:
		return user, "admin@123"
	}

	t.Fatalf("failed finding platform user/pass info")

	return "", ""
}

func prepareFunctionalDriver(
	t *testing.T,
	testName, platformName, transportName string,
) *netconf.Driver {
	host, port := getFunctionalHostIPPort(t, platformName)
	user, pass := getFunctionalHostUserPass(t, platformName)

	d, err := netconf.NewDriver(
		host,
		options.WithPort(port),
		options.WithAuthUsername(user),
		options.WithAuthPassword(pass),
		options.WithTransportType(transportName),
		options.WithAuthNoStrictKey(),
	)
	if err != nil {
		t.Fatalf("%s: encountered error creating netconf driver, error: %s", testName, err)
	}

	if platformName == platform.JuniperJunos {
		d.ForceSelfClosingTags = true
	}

	err = d.Open()
	if err != nil {
		t.Fatalf(
			"%s: encountered error opening network driver, error: %s",
			testName,
			err,
		)
	}

	return d
}

func interTestSleep() {
	if len(strings.Split(*platforms, ",")) == 1 {
		// when only running against a single platform, back to back tests tend to cause some issues
		// so basically stagger things out so the device doesn't choke. we stagger more for netconf
		// than ssh/telnet as the xr box in particular seems to not appreciate this!
		time.Sleep(1 * time.Second)

		return
	}

	if *transports == util.All {
		// when we run w/ all transports we do one transport after another, so similar to above
		// we just want to stagger things a bit.
		time.Sleep(1 * time.Second)

		return
	}
}
