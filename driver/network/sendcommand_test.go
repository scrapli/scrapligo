package network_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/scrapli/scrapligo/driver/network"

	"github.com/google/go-cmp/cmp"

	"github.com/scrapli/scrapligo/transport"

	"github.com/scrapli/scrapligo/driver/core"

	"github.com/scrapli/scrapligo/util/testhelper"
)

func platformCommandMap() map[string]string {
	return map[string]string{
		"cisco_iosxe":        "show version",
		"cisco_iosxr":        "show version",
		"cisco_nxos":         "show version",
		"arista_eos":         "show version",
		"juniper_junos":      "show version",
		"nokia_sros":         "show version",
		"nokia_sros_classic": "show version",
	}
}

func testSendCommand(
	d *network.Driver, command string,
	expectedOutput []byte, cleanFunc func(r string) string,
) func(t *testing.T) {
	return func(t *testing.T) {
		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening patched driver: %v", openErr)
		}

		r, cmdErr := d.SendCommand(command)
		if cmdErr != nil {
			t.Fatalf("failed sending command: %v", cmdErr)
		}

		if r.Failed != nil {
			t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
		}

		if diff := cmp.Diff(cleanFunc(r.Result), string(expectedOutput)); diff != "" {
			t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
		}
	}
}

func TestSendCommand(t *testing.T) {
	commandMap := platformCommandMap()

	for _, platform := range core.SupportedPlatforms() {
		command := commandMap[platform]

		sessionFile := fmt.Sprintf("../../test_data/driver/network/sendcommand/%s", platform)
		expectedFile := fmt.Sprintf(
			"../../test_data/driver/network/sendcommand/%s_expected",
			platform,
		)

		expectedOutput, expectedErr := os.ReadFile(expectedFile)
		if expectedErr != nil {
			t.Fatalf("failed opening expected output file '%s' err: %v", expectedFile, expectedErr)
		}

		d, driverErr := core.NewCoreDriver(
			"localhost",
			platform,
			testhelper.WithPatchedTransport(sessionFile),
		)

		if driverErr != nil {
			t.Fatalf("failed creating test device: %v", driverErr)
		}

		f := testSendCommand(d, command, expectedOutput, testhelper.CleanResponseNoop)

		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}

func platformFunctionalCommandMapShort() map[string]string {
	return map[string]string{
		"arista_eos": "show run | i ZTP",
	}
}

func platformFunctionalCommandExpected() map[string][]byte {
	return map[string][]byte{
		"arista_eos": []byte("logging level ZTP informational"),
	}
}

func TestFunctionalSendCommandShort(t *testing.T) {
	if !*testhelper.Functional {
		t.Skip("SKIP: functional tests skipped unless the '-functional' flag is passed")
	}

	commandMap := platformFunctionalCommandMapShort()
	expectedOutputMap := platformFunctionalCommandExpected()

	for _, transportName := range transport.SupportedTransports() {
		for platform, command := range commandMap {
			expectedOutput := expectedOutputMap[platform]

			hostConnData, ok := testhelper.FunctionalTestHosts()[platform]
			if !ok {
				t.Fatalf("no host connection data for platform type %s\n", platform)
			}

			port := hostConnData.Port
			if transportName == transport.TelnetTransportName {
				port = hostConnData.TelnetPort
			}

			d := testhelper.NewFunctionalTestDriver(
				t,
				hostConnData.Host,
				platform,
				transportName,
				port,
			)

			f := testSendCommand(d, command, expectedOutput, testhelper.CleanResponseNoop)

			t.Run(fmt.Sprintf("Platform=%s;Transport=%s", platform, transportName), f)
		}
	}
}

func platformFunctionalCommandMapLong() map[string]string {
	return map[string]string{
		"arista_eos": "show run",
	}
}

func TestFunctionalSendCommandLong(t *testing.T) {
	if !*testhelper.Functional {
		t.Skip("SKIP: functional tests skipped unless the '-functional' flag is passed")
	}

	commandMap := platformFunctionalCommandMapLong()

	for _, transportName := range transport.SupportedTransports() {
		for platform, command := range commandMap {
			cleanFunc := testhelper.CleanResponseMap()[platform]

			expectedFile := fmt.Sprintf(
				"../../test_data/driver/network/sendcommand/%s_functional_expected",
				platform,
			)

			expectedOutput, expectedErr := os.ReadFile(expectedFile)
			if expectedErr != nil {
				t.Fatalf(
					"failed opening expected output file '%s' err: %v",
					expectedFile,
					expectedErr,
				)
			}

			hostConnData, ok := testhelper.FunctionalTestHosts()[platform]
			if !ok {
				t.Fatalf("no host connection data for platform type %s\n", platform)
			}

			port := hostConnData.Port
			if transportName == transport.TelnetTransportName {
				port = hostConnData.TelnetPort
			}

			d := testhelper.NewFunctionalTestDriver(
				t,
				hostConnData.Host,
				platform,
				transportName,
				port,
			)

			f := testSendCommand(d, command, expectedOutput, cleanFunc)

			t.Run(fmt.Sprintf("Platform=%s;Transport=%s", platform, transportName), f)
		}
	}
}
