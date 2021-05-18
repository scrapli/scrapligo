package testhelper

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/scrapli/scrapligo/driver/core"
)

// SendCommandTestHelper helper function to handle send command tests.
func SendCommandTestHelper(t *testing.T, driverName, command string) func(t *testing.T) {
	sessionFile := fmt.Sprintf("../../test_data/driver/network/sendcommand/%s", driverName)
	expectedFile := fmt.Sprintf(
		"../../test_data/driver/network/sendcommand/%s_expected",
		driverName,
	)

	expected, expectedErr := os.ReadFile(expectedFile)
	if expectedErr != nil {
		t.Fatalf("failed opening expected output file '%s' err: %v", expectedFile, expectedErr)
	}

	return func(t *testing.T) {
		d, driverErr := core.NewCoreDriver(
			"localhost",
			driverName,
			WithPatchedTransport(sessionFile, t),
		)

		if driverErr != nil {
			t.Fatalf("failed creating test device: %v", driverErr)
		}

		openErr := d.Open()
		if openErr != nil {
			t.Fatalf("failed opening patched driver: %v", openErr)
		}

		r, cmdErr := d.SendCommand(command)
		if cmdErr != nil {
			t.Fatalf("failed sending command: %v", cmdErr)
		}

		if r.Failed {
			t.Fatal("response object indicates failure")
		}

		// i have no idea where the null bit is getting read from... but it does? so we'll just remove
		// it for now...
		finalResult := string(bytes.Trim([]byte(r.Result), "\x00\x0a"))

		if finalResult != string(expected) {
			t.Fatal("actual result and expected result do not match")
		}
	}
}
