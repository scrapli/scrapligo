package network_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/scrapli/scrapligo/driver/network"

	"github.com/google/go-cmp/cmp"

	"github.com/scrapli/scrapligo/transport"

	"github.com/scrapli/scrapligo/util/testhelper"
)

func testSendCommand(
	d *network.Driver, command string,
	expectedOutput []byte, cleanFunc func(r string) string,
) func(t *testing.T) {
	return func(t *testing.T) {
		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening driver: %v", openErr)
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

	for platform, command := range commandMap {
		sessionFile := fmt.Sprintf(
			"../../test_data/driver/network/sendcommand/%s_session_short",
			platform,
		)

		expectedFile := fmt.Sprintf(
			"../../test_data/driver/network/expected/%s_short_expected",
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

		d := testhelper.CreatePatchedDriver(t, sessionFile, platform)

		f := testSendCommand(
			d,
			command,
			expectedOutput,
			testhelper.GetCleanFunc(platform),
		)

		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}

func testFunctionalSendCommandCommon(
	t *testing.T,
	command, expectedFile, platform, transportName string,
) {
	if !testhelper.RunPlatform(platform) {
		t.Logf("skip; platform %s deselected for testing\n", platform)
		return
	}

	hostConnData, ok := functionalTestHosts()[platform]
	if !ok {
		t.Logf("skip; no host connection data for platform type %s\n", platform)
		return
	}

	expectedOutput, expectedErr := os.ReadFile(expectedFile)
	if expectedErr != nil {
		t.Fatalf(
			"failed opening expected output file '%s' err: %v",
			expectedFile,
			expectedErr,
		)
	}

	port := hostConnData.Port
	if transportName == transport.TelnetTransportName {
		port = hostConnData.TelnetPort
	}

	d := newFunctionalTestDriver(
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
		testhelper.GetCleanFunc(platform),
	)

	t.Run(fmt.Sprintf("Platform=%s;Transport=%s", platform, transportName), f)
}

func TestFunctionalSendCommandShort(t *testing.T) {
	if !*testhelper.Functional {
		t.Skip("skip: functional tests skipped unless the '-functional' flag is passed")
	}

	commandMap := platformCommandMapShort()

	for _, transportName := range transport.SupportedTransports() {
		if !testhelper.RunTransport(transportName) {
			t.Logf("skip; transport %s deselected for testing\n", transportName)
			continue
		}

		for platform, command := range commandMap {
			expectedFile := fmt.Sprintf(
				"../../test_data/driver/network/expected/%s_short_expected",
				platform,
			)

			testFunctionalSendCommandCommon(t, command, expectedFile, platform, transportName)
		}
	}
}

func TestFunctionalSendCommandLong(t *testing.T) {
	if !*testhelper.Functional {
		t.Skip("skip: functional tests skipped unless the '-functional' flag is passed")
	}

	commandMap := platformCommandMapLong()

	for _, transportName := range transport.SupportedTransports() {
		if !testhelper.RunTransport(transportName) {
			t.Logf("skip; transport %s deselected for testing\n", transportName)
			continue
		}

		for platform, command := range commandMap {
			expectedFile := fmt.Sprintf(
				"../../test_data/driver/network/expected/%s_long_expected",
				platform,
			)

			testFunctionalSendCommandCommon(t, command, expectedFile, platform, transportName)
		}
	}
}
