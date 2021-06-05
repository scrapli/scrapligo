package channel_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util/testhelper"
)

func TestSendInput(t *testing.T) {
	fakeSession := "sendinput"

	expectedFile := "../test_data/channel/sendinput_expected"

	expected, expectedErr := os.ReadFile(expectedFile)
	if expectedErr != nil {
		t.Fatalf("failed opening expected output file '%s' err: %v", expectedFile, expectedErr)
	}

	c := testhelper.NewPatchedChannel(t, &fakeSession)

	actual, promptErr := c.SendInput("show version", false, false, 0)

	if promptErr != nil {
		t.Fatalf("error sending input to mock channel: %v", promptErr)
	}

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}
