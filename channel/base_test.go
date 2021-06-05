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

	testTransp := testhelper.FetchTestTransport(c, t)

	if diff := cmp.Diff(testTransp.CapturedWrites[0], channelInput); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}

func TestSendReturn(t *testing.T) {
	c := testhelper.NewPatchedChannel()

	writeErr := c.SendReturn()

	if writeErr != nil {
		t.Fatalf("error writing to mock channel: %v", writeErr)
	}

	testTransp := testhelper.FetchTestTransport(c, t)

	if diff := cmp.Diff(testTransp.CapturedWrites[0], []byte(*c.CommsReturnChar)); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}

func TestWriteAndReturn(t *testing.T) {
	c := testhelper.NewPatchedChannel()

	channelInput := []byte("something witty")
	finalChannelInput := [][]byte{channelInput, []byte(*c.CommsReturnChar)}

	writeErr := c.WriteAndReturn(channelInput, false)

	if writeErr != nil {
		t.Fatalf("error writing to mock channel: %v", writeErr)
	}

	testTransp := testhelper.FetchTestTransport(c, t)

	if diff := cmp.Diff(testTransp.CapturedWrites, finalChannelInput); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}
