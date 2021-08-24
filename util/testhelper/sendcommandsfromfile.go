package testhelper

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/scrapli/scrapligo/driver/core"
)

// SendCommandsFromFileTestHelper helper function to handle send commands from file tests.
func SendCommandsFromFileTestHelper(t *testing.T, driverName string) func(t *testing.T) {
	sessionFile := fmt.Sprintf("../../test_data/driver/network/sendcommandsfromfile/%s", driverName)
	expectedFileOne := fmt.Sprintf(
		"../../test_data/driver/network/sendcommandsfromfile/%s_expected_one",
		driverName,
	)
	expectedFileTwo := fmt.Sprintf(
		"../../test_data/driver/network/sendcommandsfromfile/%s_expected_two",
		driverName,
	)

	expectedOne, expectedErr := os.ReadFile(expectedFileOne)
	if expectedErr != nil {
		t.Fatalf("failed opening expected output file '%s' err: %v", expectedFileOne, expectedErr)
	}

	expectedTwo, expectedErr := os.ReadFile(expectedFileTwo)
	if expectedErr != nil {
		t.Fatalf("failed opening expected output file '%s' err: %v", expectedFileOne, expectedErr)
	}

	return func(t *testing.T) {
		d, driverErr := core.NewCoreDriver(
			"localhost",
			driverName,
			WithPatchedTransport(sessionFile),
		)

		if driverErr != nil {
			t.Fatalf("failed creating test device: %v", driverErr)
		}

		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening patched driver: %v", openErr)
		}

		r, cmdErr := d.SendCommandsFromFile(
			fmt.Sprintf(
				"../../test_data/driver/network/sendcommandsfromfile/%s_commands",
				driverName,
			),
		)
		if cmdErr != nil {
			t.Fatalf("failed sending command: %v", cmdErr)
		}

		if r.Failed != nil {
			t.Fatalf("response object indicates failure; error: %+v\n", r.Failed)
		}

		// i have no idea where the null bit is getting read from... but it does? so we'll just remove
		// it for now...
		responseOne := r.Responses[0].Result
		responseTwo := r.Responses[1].Result

		finalResultOne := string(bytes.Trim([]byte(responseOne), "\x00\x0a"))
		finalResultTwo := string(bytes.Trim([]byte(responseTwo), "\x00\x0a"))

		if finalResultOne != string(expectedOne) {
			t.Fatal(
				"actual result one and expected result do not match",
			)
		}

		if finalResultTwo != string(expectedTwo) {
			t.Fatal("actual result two and expected result do not match")
		}
	}
}
