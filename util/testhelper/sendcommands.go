package testhelper

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/driver/core"
)

// SendCommandsTestHelper helper function to handle send commands tests.
func SendCommandsTestHelper(t *testing.T, driverName string, commands []string) func(t *testing.T) {
	sessionFile := fmt.Sprintf("../../test_data/driver/network/sendcommands/%s", driverName)
	expectedFileOne := fmt.Sprintf(
		"../../test_data/driver/network/sendcommands/%s_expected_one",
		driverName,
	)
	expectedFileTwo := fmt.Sprintf(
		"../../test_data/driver/network/sendcommands/%s_expected_two",
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

		r, cmdErr := d.SendCommands(commands)
		if cmdErr != nil {
			t.Fatalf("failed sending command: %v", cmdErr)
		}

		if r.Failed() {
			t.Fatal("response object indicates failure")
		}

		// i have no idea where the null bit is getting read from... but it does? so we'll just remove
		// it for now...
		responseOne := r.Responses[0].Result
		responseTwo := r.Responses[1].Result

		finalResultOne := string(bytes.Trim([]byte(responseOne), "\x00\x0a"))
		finalResultTwo := string(bytes.Trim([]byte(responseTwo), "\x00\x0a"))

		if diff := cmp.Diff(finalResultOne, string(expectedOne)); diff != "" {
			t.Errorf("actual result one and expected result do not match (-want +got):\n%s", diff)
		}

		if diff := cmp.Diff(finalResultTwo, string(expectedTwo)); diff != "" {
			t.Errorf("actual result two and expected result do not match (-want +got):\n%s", diff)
		}
	}
}
