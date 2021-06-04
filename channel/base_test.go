package channel_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/scrapli/scrapligo/util/testhelper"
)

func TestWrite(t *testing.T) {
	c := testhelper.NewPatchedChannel()

	channelInput := []byte("something witty")

	writeErr := c.Write(channelInput, false)

	if writeErr != nil {
		t.Fatalf("error writing to mock channel: %v", writeErr)
	}

	capturedWrites := testhelper.FetchCapturedWrites(c.Transport, t)

	if diff := cmp.Diff(capturedWrites[0], channelInput); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}
