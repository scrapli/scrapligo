package channel_test

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/scrapli/scrapligo/util/testhelper"
)

func TestWrite(t *testing.T) {
	c := testhelper.NewPatchedChannel(t, nil)

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
	c := testhelper.NewPatchedChannel(t, nil)

	writeErr := c.SendReturn()

	if writeErr != nil {
		t.Fatalf("error writing to mock channel: %v", writeErr)
	}

	testTransp := testhelper.FetchTestTransport(c, t)

	if diff := cmp.Diff(testTransp.CapturedWrites[0], []byte(c.CommsReturnChar)); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}

func TestWriteAndReturn(t *testing.T) {
	c := testhelper.NewPatchedChannel(t, nil)

	channelInput := []byte("something witty")
	finalChannelInput := [][]byte{channelInput, []byte(c.CommsReturnChar)}

	writeErr := c.WriteAndReturn(channelInput, false)

	if writeErr != nil {
		t.Fatalf("error writing to mock channel: %v", writeErr)
	}

	testTransp := testhelper.FetchTestTransport(c, t)

	if diff := cmp.Diff(testTransp.CapturedWrites, finalChannelInput); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}

func TestRead(t *testing.T) {
	fakeSession := "read"
	expectedFile := "../test_data/channel/read_expected"

	expected, expectedErr := os.ReadFile(expectedFile)
	if expectedErr != nil {
		t.Fatalf("failed opening expected output file '%s' err: %v", expectedFile, expectedErr)
	}

	c := testhelper.NewPatchedChannel(t, &fakeSession)
	testhelper.SetTestTransportStandardReadSize(c, t)

	b, readErr := c.Read()

	if readErr != nil {
		t.Fatalf("error reading from mock channel: %v", readErr)
	}

	if diff := cmp.Diff(b, expected); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}

func TestRestructureOutput(t *testing.T) {
	output := []byte("   some output\nlocalhost#   ")
	expected := []byte("   some output\nlocalhost#")

	c := testhelper.NewPatchedChannel(t, nil)

	actual := c.RestructureOutput(output, false)

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}

func TestRestructureOutputStripPrompt(t *testing.T) {
	output := []byte("   some output   \nlocalhost# ")
	expected := []byte("   some output")

	c := testhelper.NewPatchedChannel(t, nil)

	actual := c.RestructureOutput(output, true)

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}

func TestDetermineOperationTimeoutDefault(t *testing.T) {
	expected := 30 * time.Second

	c := testhelper.NewPatchedChannel(t, nil)

	actual := c.DetermineOperationTimeout(0)

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}

func TestDetermineOperationTimeoutMax(t *testing.T) {
	expected := 24 * time.Hour

	c := testhelper.NewPatchedChannel(t, nil)
	c.TimeoutOps = 0 * time.Second

	actual := c.DetermineOperationTimeout(0)

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}

func TestDetermineOperationTimeoutPerCommand(t *testing.T) {
	expected := 1 * time.Minute

	c := testhelper.NewPatchedChannel(t, nil)
	c.TimeoutOps = 0 * time.Second
	actual := c.DetermineOperationTimeout(1 * time.Minute)

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}

func TestFormatLogMessage(t *testing.T) {
	expected := "debug::localhost::22::some log message"

	c := testhelper.NewPatchedChannel(t, nil)

	actual := c.FormatLogMessage("debug", "some log message")

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}
