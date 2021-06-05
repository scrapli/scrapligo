package channel_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util/testhelper"
)

func TestGetPrompt(t *testing.T) {
	expected := "localhost#"

	fakeSession := "getprompt"

	c := testhelper.NewPatchedChannel(t, &fakeSession)

	actual, promptErr := c.GetPrompt()

	if promptErr != nil {
		t.Fatalf("error getting prompt from mock channel: %v", promptErr)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}
