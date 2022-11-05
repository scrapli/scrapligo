package network_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/platform"
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
) (*network.Driver, *transport.File) {
	d, err := network.NewDriver(
		"dummy",
		options.WithTransportType(transport.FileTransport),
		options.WithFileTransportFile(resolveFile(t, payloadFile)),
		options.WithTransportReadSize(1),
		options.WithReadDelay(0),
		options.WithDefaultDesiredPriv("privilege-exec"),
		options.WithPrivilegeLevels(map[string]*network.PrivilegeLevel{
			"exec": {
				Pattern:        `(?im)^[\w.\-@/:]{1,63}>$`,
				Name:           "exec",
				PreviousPriv:   "",
				Deescalate:     "",
				Escalate:       "",
				EscalateAuth:   false,
				EscalatePrompt: "",
			},
			"privilege-exec": {
				Pattern:        `(?im)^[\w.\-@/:]{1,63}#$`,
				Name:           "privilege-exec",
				PreviousPriv:   "exec",
				Deescalate:     "disable",
				Escalate:       "enable",
				EscalateAuth:   true,
				EscalatePrompt: `(?im)^(?:enable\s){0,1}password:\s?$`,
			},
			"configuration": {
				Pattern:        `(?im)^[\w.\-@/:]{1,63}\([\w.\-@/:+]{0,32}\)#$`,
				NotContains:    []string{"tcl)"},
				Name:           "configuration",
				PreviousPriv:   "privilege-exec",
				Deescalate:     "end",
				Escalate:       "configure terminal",
				EscalateAuth:   false,
				EscalatePrompt: "",
			},
		}),
	)
	if err != nil {
		t.Fatalf("%s: encountered error creating network Driver, error: %s", testName, err)
	}

	err = d.Channel.Open()
	if err != nil {
		t.Fatalf("%s: encountered error opening Channel, error: %s", testName, err)
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

func getFunctionalHostIPPortLinuxOrRemote(
	t *testing.T,
	platformName string,
) (host string, port int) {
	switch platformName {
	case platform.CiscoIosxe:
		host = util.GetEnvStrOrDefault("SCRAPLIGO_CISCO_IOSXE_HOST", "172.20.20.11")

		return host, 22
	case platform.CiscoIosxr:
		host = util.GetEnvStrOrDefault("SCRAPLIGO_CISCO_IOSXR_HOST", "172.20.20.12")

		return host, 22
	case platform.CiscoNxos:
		host = util.GetEnvStrOrDefault("SCRAPLIGO_CISCO_NXOS_HOST", "172.20.20.13")

		return host, 22
	case platform.AristaEos:
		host = util.GetEnvStrOrDefault("SCRAPLIGO_ARISTA_EOS_HOST", "172.20.20.14")

		return host, 22
	case platform.JuniperJunos:
		host = util.GetEnvStrOrDefault("SCRAPLIGO_JUNIPER_JUNOS_HOST", "172.20.20.15")

		return host, 22
	case platform.NokiaSrl:
		host = util.GetEnvStrOrDefault("SCRAPLIGO_NOKIA_SRL_HOST", "172.20.20.16")

		return host, 22
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
		return host, 21022
	case platform.CiscoIosxr:
		return host, 22022
	case platform.CiscoNxos:
		return host, 23022
	case platform.AristaEos:
		return host, 24022
	case platform.JuniperJunos:
		return host, 25022
	case platform.NokiaSrl:
		return host, 26022
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
	case platform.NokiaSrl:
		return user, pass
	}

	t.Fatalf("failed finding platform user/pass info")

	return "", ""
}

func prepareFunctionalDriver(
	t *testing.T,
	testName, platformName, transportName string,
) *network.Driver {
	host, port := getFunctionalHostIPPort(t, platformName)
	user, pass := getFunctionalHostUserPass(t, platformName)

	if transportName == transport.TelnetTransport {
		if platformName != platform.CiscoIosxe {
			t.Skip("only testing telnet on iosxe at the moment")
		}

		port++
	}

	p, err := platform.NewPlatform(
		fmt.Sprintf("%s.yaml", platformName),
		host,
		options.WithPort(port),
		options.WithAuthUsername(user),
		options.WithAuthPassword(pass),
		options.WithTransportType(transportName),
		options.WithAuthNoStrictKey(),
		// obviously only relevant for system transport, but will be ignored for others
		// also should only be necessary for iosxe.
		options.WithSystemTransportOpenArgs(
			[]string{
				"-o",
				"KexAlgorithms=+diffie-hellman-group-exchange-sha1,diffie-hellman-group14-sha1",
				"-o",
				// note that PubkeyAcceptedKeyTypes works on older versions of openssh, whereas
				// PubKeyAcceptedAlgorithms is the option on >=8.5, runners in actions use older
				// version so we'll roll with th older type here.
				"PubkeyAcceptedKeyTypes=+ssh-rsa",
				"-o",
				"HostKeyAlgorithms=+ssh-dss,ssh-rsa,rsa-sha2-512,rsa-sha2-256,ssh-rsa,ssh-ed25519",
			},
		),
	)
	if err != nil {
		t.Fatalf("%s: encountered error creating platform, error: %s", testName, err)
	}

	d, err := p.GetNetworkDriver()
	if err != nil {
		t.Fatalf(
			"%s: encountered error fetching network driver from platform, error: %s",
			testName,
			err,
		)
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
		// so basically stagger things out so the device doesn't choke.
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
