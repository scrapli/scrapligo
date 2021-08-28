package network_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/transport"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/util/testhelper"
)

func compareResults(
	t *testing.T,
	r *base.MultiResponse,
	expectedOutput [][]byte,
	cleanFunc func(r string) string,
) {
	if r.Failed != nil {
		t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
	}

	rOne := r.Responses[0].Result
	rTwo := r.Responses[1].Result

	if diff := cmp.Diff(cleanFunc(rOne), string(expectedOutput[0])); diff != "" {
		t.Errorf("actual result one and expected result do not match (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(cleanFunc(rTwo), string(expectedOutput[1])); diff != "" {
		t.Errorf("actual result two and expected result do not match (-want +got):\n%s", diff)
	}
}

func testSendCommands(
	d *network.Driver, commands []string,
	expectedOutput [][]byte, cleanFunc func(r string) string,
) func(t *testing.T) {
	return func(t *testing.T) {
		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening driver: %v", openErr)
		}

		r, cmdErr := d.SendCommands(commands)
		if cmdErr != nil {
			t.Fatalf("failed sending commands: %v", cmdErr)
		}

		compareResults(t, r, expectedOutput, cleanFunc)
	}
}

func testSendCommandsFromFile(
	d *network.Driver,
	commandsFile string,
	expectedOutput [][]byte,
	cleanFunc func(r string) string,
) func(t *testing.T) {
	return func(t *testing.T) {
		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening driver: %v", openErr)
		}

		r, cmdErr := d.SendCommandsFromFile(commandsFile)
		if cmdErr != nil {
			t.Fatalf("failed sending commands: %v", cmdErr)
		}

		compareResults(t, r, expectedOutput, cleanFunc)
	}
}

func testSendCommandsCommon(
	t *testing.T,
	platform string,
) (d *network.Driver, expectedOutputs [][]byte) {
	sessionFile := fmt.Sprintf(
		"../../test_data/driver/network/sendcommands/%s_session",
		platform,
	)

	expectedFileOne := fmt.Sprintf(
		"../../test_data/driver/network/expected/%s_short_expected",
		platform,
	)

	expectedFileTwo := fmt.Sprintf(
		"../../test_data/driver/network/expected/%s_long_expected",
		platform,
	)

	expectedOutputOne, expectedErr := os.ReadFile(expectedFileOne)
	if expectedErr != nil {
		t.Fatalf(
			"failed opening expected output file '%s' err: %v",
			expectedFileOne,
			expectedErr,
		)
	}

	expectedOutputTwo, expectedErr := os.ReadFile(expectedFileTwo)
	if expectedErr != nil {
		t.Fatalf(
			"failed opening expected output file '%s' err: %v",
			expectedFileTwo,
			expectedErr,
		)
	}

	d = testhelper.CreatePatchedDriver(t, sessionFile, platform)

	return d, [][]byte{expectedOutputOne, expectedOutputTwo}
}

func TestSendCommands(t *testing.T) {
	commandMapShort := platformCommandMapShort()
	commandMapLong := platformCommandMapLong()

	for _, platform := range core.SupportedPlatforms() {
		commandOne := commandMapShort[platform]
		commandTwo := commandMapLong[platform]

		d, expectedOutputs := testSendCommandsCommon(t, platform)

		f := testSendCommands(
			d,
			[]string{commandOne, commandTwo},
			expectedOutputs,
			testhelper.GetCleanFunc(platform),
		)

		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}

func TestSendCommandsFromFile(t *testing.T) {
	for _, platform := range core.SupportedPlatforms() {
		d, expectedOutputs := testSendCommandsCommon(t, platform)

		f := testSendCommandsFromFile(
			d,
			fmt.Sprintf(
				"../../test_data/driver/network/sendcommandsfromfile/%s_commands",
				platform,
			),
			expectedOutputs,
			testhelper.GetCleanFunc(platform),
		)

		t.Run(fmt.Sprintf("Platform=%s", platform), f)
	}
}

func testFunctionalSendCommandsCommon(
	t *testing.T,
	platform, transportName string,
) (d *network.Driver, expectedOutputs [][]byte) {
	if !testhelper.RunPlatform(platform) {
		t.Logf("skip; platform %s deselected for testing\n", platform)
		return
	}

	hostConnData, ok := functionalTestHosts()[platform]
	if !ok {
		t.Logf("skip; no host connection data for platform type %s\n", platform)
		return
	}

	expectedFileOne := fmt.Sprintf(
		"../../test_data/driver/network/expected/%s_short_expected",
		platform,
	)

	expectedFileTwo := fmt.Sprintf(
		"../../test_data/driver/network/expected/%s_long_expected",
		platform,
	)

	expectedOutputOne, expectedErr := os.ReadFile(expectedFileOne)
	if expectedErr != nil {
		t.Fatalf(
			"failed opening expected output file '%s' err: %v",
			expectedFileOne,
			expectedErr,
		)
	}

	expectedOutputTwo, expectedErr := os.ReadFile(expectedFileTwo)
	if expectedErr != nil {
		t.Fatalf(
			"failed opening expected output file '%s' err: %v",
			expectedFileTwo,
			expectedErr,
		)
	}

	port := hostConnData.Port
	if transportName == transport.TelnetTransportName {
		port = hostConnData.TelnetPort
	}

	d = newFunctionalTestDriver(
		t,
		hostConnData.Host,
		platform,
		transportName,
		port,
	)

	return d, [][]byte{expectedOutputOne, expectedOutputTwo}
}

func TestFunctionalSendCommands(t *testing.T) {
	if !*testhelper.Functional {
		t.Skip("skip: functional tests skipped unless the '-functional' flag is passed")
	}

	commandMapShort := platformCommandMapShort()
	commandMapLong := platformCommandMapLong()

	for _, transportName := range transport.SupportedTransports() {
		if !testhelper.RunTransport(transportName) {
			t.Logf("skip; transport %s deselected for testing\n", transportName)
			continue
		}

		for _, platform := range core.SupportedPlatforms() {
			if !testhelper.RunPlatform(platform) {
				t.Logf("skip; platform %s deselected for testing\n", platform)
				continue
			}

			commandOne := commandMapShort[platform]
			commandTwo := commandMapLong[platform]

			d, expectedOutputs := testFunctionalSendCommandsCommon(t, platform, transportName)
			if d == nil {
				// no connection data or some reason to skip
				continue
			}

			f := testSendCommands(
				d,
				[]string{commandOne, commandTwo},
				expectedOutputs,
				testhelper.GetCleanFunc(platform),
			)

			t.Run(fmt.Sprintf("Platform=%s;Transport=%s", platform, transportName), f)
		}
	}
}

func TestFunctionalSendCommandsFromFile(t *testing.T) {
	if !*testhelper.Functional {
		t.Skip("skip: functional tests skipped unless the '-functional' flag is passed")
	}

	for _, transportName := range transport.SupportedTransports() {
		if !testhelper.RunTransport(transportName) {
			t.Logf("skip; transport %s deselected for testing\n", transportName)
			continue
		}

		for _, platform := range core.SupportedPlatforms() {
			if !testhelper.RunPlatform(platform) {
				t.Logf("skip; platform %s deselected for testing\n", platform)
				continue
			}

			d, expectedOutputs := testFunctionalSendCommandsCommon(t, platform, transportName)
			if d == nil {
				// no connection data or some reason to skip
				continue
			}

			f := testSendCommandsFromFile(
				d,
				fmt.Sprintf(
					"../../test_data/driver/network/sendcommandsfromfile/%s_commands",
					platform,
				),
				expectedOutputs,
				testhelper.GetCleanFunc(platform),
			)

			t.Run(fmt.Sprintf("Platform=%s;Transport=%s", platform, transportName), f)
		}
	}
}
