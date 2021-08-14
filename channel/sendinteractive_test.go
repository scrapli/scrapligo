package channel_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/util/testhelper"
)

func TestSendInteractive(t *testing.T) {
	fakeSession := "sendinteractive"

	expectedFile := "../test_data/channel/sendinteractive"

	expected, expectedErr := os.ReadFile(expectedFile)
	if expectedErr != nil {
		t.Fatalf("failed opening expected output file '%s' err: %v", expectedFile, expectedErr)
	}

	c := testhelper.NewPatchedChannel(t, &fakeSession)

	events := []*channel.SendInteractiveEvent{
		{
			ChannelInput:    "clear logging",
			ChannelResponse: "[confirm]",
			HideInput:       false,
		},
		{
			ChannelInput:    "",
			ChannelResponse: "",
			HideInput:       false,
		},
	}
	actual, promptErr := c.SendInteractive(events, nil, 0)

	if promptErr != nil {
		t.Fatalf("error sending input to mock channel: %v", promptErr)
	}

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}
