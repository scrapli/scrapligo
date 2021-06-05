package testhelper

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/scrapli/scrapligo/channel"
)

// NewPatchedChannel create a new channel that is patched with testing transport.
func NewPatchedChannel(t *testing.T, sessionFile *string) *channel.Channel {
	transport := &TestingTransport{}

	if sessionFile != nil {
		finalSessionFile := fmt.Sprintf("../test_data/channel/%s", *sessionFile)

		f, err := os.Open(finalSessionFile)
		if err != nil {
			t.Fatalf("failed opening transport session file '%s' err: %v", finalSessionFile, err)
		}

		transport.FakeSession = f
	}

	returnChar := "\n"
	commsPromptPattern := regexp.MustCompile(`(?im)^localhost#\s?$`)
	timeoutOps := 30 * time.Second

	c := &channel.Channel{
		Transport:              transport,
		CommsPromptPattern:     commsPromptPattern,
		CommsReturnChar:        &returnChar,
		CommsPromptSearchDepth: 255,
		TimeoutOps:             &timeoutOps,
		Host:                   "localhost",
		Port:                   22,
	}

	return c
}

// FetchTestTransport fetch the TestTransport object so we can do operations on attributes that only
// the test transport has.
func FetchTestTransport(c *channel.Channel, t *testing.T) *TestingTransport {
	testTransp, ok := c.Transport.(*TestingTransport)

	if !ok {
		t.Fatalf("this should not happen; TestingTransport patching failed somehow :(")
	}

	return testTransp
}

// SetTestTransportStandardReadSize set the TestTransport read size to the "normal" value of 65535 -
// this is living in the channel file as its only necessary to modify for channel test operations.
func SetTestTransportStandardReadSize(c *channel.Channel, t *testing.T) {
	testTransp, ok := c.Transport.(*TestingTransport)

	if !ok {
		t.Fatalf("this should not happen; TestingTransport patching failed somehow :(")
	}

	readSize := 65535
	testTransp.ReadSize = &readSize
}
