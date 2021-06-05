package testhelper

import (
	"testing"

	"github.com/scrapli/scrapligo/channel"
)

// NewPatchedChannel create a new channel that is patched with testing transport.
func NewPatchedChannel() *channel.Channel {
	transport := &TestingTransport{}

	returnChar := "\n"

	c := &channel.Channel{
		Transport:              transport,
		CommsPromptPattern:     nil,
		CommsReturnChar:        &returnChar,
		CommsPromptSearchDepth: 255,
		TimeoutOps:             nil,
		Host:                   "localhost",
		Port:                   22,
	}

	return c
}

func FetchTestTransport(c *channel.Channel, t *testing.T) *TestingTransport {
	testTransp, ok := c.Transport.(*TestingTransport)

	if !ok {
		t.Fatalf("this should not happen; TestingTransport patching failed somehow :(")
	}

	return testTransp
}
