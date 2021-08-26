package network_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/scrapli/scrapligo/driver/core"

	"github.com/scrapli/scrapligo/driver/network"

	"github.com/google/go-cmp/cmp"

	"github.com/scrapli/scrapligo/transport"

	"github.com/scrapli/scrapligo/util/testhelper"
)

func platformCommandMapShort() map[string]string {
	return map[string]string{
		"cisco_iosxe":        "show run | i hostname",
		"cisco_iosxr":        "show run | i MgmtEth0",
		"cisco_nxos":         "show run | i scp-server",
		"arista_eos":         "show run | i ZTP",
		"juniper_junos":      "show configuration | match 10.0.0.15",
		"nokia_sros":         "show version",
		"nokia_sros_classic": "show version",
	}
}

func platformCommandExpected() map[string][]byte {
	return map[string][]byte{
		"cisco_iosxe": []byte("hostname csr1000v"),
		"cisco_iosxr": []byte(
			"TIME_STAMP_REPLACED\nBuilding configuration...\ninterface MgmtEth0/RP0/CPU0/0",
		),
		"cisco_nxos":    []byte("feature scp-server"),
		"arista_eos":    []byte("logging level ZTP informational"),
		"juniper_junos": []byte("                address 10.0.0.15/24;"),
		"nokia_sros": []byte(
			"TiMOS-B-20.10.R3 both/x86_64 Nokia 7750 SR Copyright (c) 2000-2021 Nokia.\nAll " +
				"rights reserved. All use subject to applicable license agreements.\nBuilt on " +
				"Wed Jan 27 13:21:10 PST 2021 by builder in /builds/c/2010B/R3/panos/main/sros",
		),
		"nokia_sros_classic": []byte(
			"TiMOS-B-20.10.R3 both/x86_64 Nokia 7750 SR Copyright (c) 2000-2021 Nokia.\nAll " +
				"rights reserved. All use subject to applicable license agreements.\nBuilt on " +
				"Wed Jan 27 13:21:10 PST 2021 by builder in /builds/c/2010B/R3/panos/main/sros",
		),
	}
}

func platformCommandMapLong() map[string]string {
	return map[string]string{
		"arista_eos":    "show run",
		"juniper_junos": "show configuration",
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
	commandMap := platformCommandMapShort()
	expectedOutputMap := platformCommandExpected()

	for platform, command := range commandMap {
		sessionFile := fmt.Sprintf("../../test_data/driver/network/sendcommand/%s", platform)

		expectedOutput := expectedOutputMap[platform]

		d, driverErr := core.NewCoreDriver(
			"localhost",
			platform,
			testhelper.WithPatchedTransport(sessionFile),
		)

		if driverErr != nil {
			t.Fatalf("failed creating test device: %v", driverErr)
		}

		f := testSendCommand(
			d,
			command,
			expectedOutput,
			testhelper.GetCleanFunc(platform, expectedOutput),
		)

		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}

func TestFunctionalSendCommandShort(t *testing.T) {
	if !*testhelper.Functional {
		t.Skip("SKIP: functional tests skipped unless the '-functional' flag is passed")
	}

	commandMap := platformCommandMapShort()
	expectedOutputMap := platformCommandExpected()

	for _, transportName := range transport.SupportedTransports() {
		for platform, command := range commandMap {
			expectedOutput := expectedOutputMap[platform]

			hostConnData, ok := testhelper.FunctionalTestHosts()[platform]
			if !ok {
				t.Logf("skip; no host connection data for platform type %s\n", platform)
				continue
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

			f := testSendCommand(
				d,
				command,
				expectedOutput,
				testhelper.GetCleanFunc(platform, expectedOutput),
			)

			t.Run(fmt.Sprintf("Platform=%s;Transport=%s", platform, transportName), f)
		}
	}
}

func TestFunctionalSendCommandLong(t *testing.T) {
	if !*testhelper.Functional {
		t.Skip("SKIP: functional tests skipped unless the '-functional' flag is passed")
	}

	commandMap := platformCommandMapLong()

	for _, transportName := range transport.SupportedTransports() {
		for platform, command := range commandMap {
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
				t.Logf("skip; no host connection data for platform type %s\n", platform)
				continue
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

			f := testSendCommand(
				d,
				command,
				expectedOutput,
				testhelper.GetCleanFunc(platform, expectedOutput),
			)

			t.Run(fmt.Sprintf("Platform=%s;Transport=%s", platform, transportName), f)
		}
	}
}
