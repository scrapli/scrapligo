package netconf

import (
	"fmt"
	"strconv"

	"github.com/scrapli/scrapligo/util"
)

const (
	numSessionIDMatches = 2
)

// ServerHasCapability returns true if the server supports capability s, otherwise false.
func (d *Driver) ServerHasCapability(s string) bool {
	for _, serverCapability := range d.serverCapabilities {
		if serverCapability == s {
			return true
		}
	}

	return false
}

// ServerCapabilities returns the list of capabilities the server
// sent in the initial Hello message.
func (d *Driver) ServerCapabilities() []string {
	caps := make([]string, 0, len(d.serverCapabilities))

	return append(caps, d.serverCapabilities...)
}

// SessionID returns the session ID sent by the server in the initial Hello message.
func (d *Driver) SessionID() uint64 {
	return d.sessionID
}

func (d *Driver) processServerCapabilities() error {
	b, err := d.Channel.ReadUntilPrompt()
	if err != nil {
		return err
	}

	ncPatterns := getNetconfPatterns()

	serverHelloMatch := ncPatterns.hello.Match(b)

	if !serverHelloMatch {
		return fmt.Errorf("%w: did not find server hello", util.ErrNetconfError)
	}

	// rather than deal w/ xml like scrapli python does, just regex the caps out
	serverCapabilitiesMatches := ncPatterns.capability.FindAllSubmatch(b, -1)

	d.serverCapabilities = make([]string, 0, len(serverCapabilitiesMatches))
	for _, match := range serverCapabilitiesMatches {
		d.serverCapabilities = append(d.serverCapabilities, string(match[1]))
	}

	// extract session id if it exists in the hello message
	sessionIDMatch := ncPatterns.sessionID.FindSubmatch(b)
	if len(sessionIDMatch) != numSessionIDMatches {
		return nil
	}

	i, err := strconv.Atoi(string(sessionIDMatch[1]))
	if err != nil {
		return fmt.Errorf("%w: invalid sessionID=%s: %v", util.ErrNetconfError,
			string(sessionIDMatch[1]), err)
	}

	d.sessionID = uint64(i)

	return nil
}

func (d *Driver) determineVersion() error {
	if d.ServerHasCapability(v1Dot1Cap) { //nolint: gocritic
		d.SelectedVersion = V1Dot1
	} else if d.ServerHasCapability(v1Dot0Cap) {
		d.SelectedVersion = V1Dot0
	} else {
		return fmt.Errorf("%w: capabilities exchange failed", util.ErrNetconfError)
	}

	switch d.PreferredVersion {
	case V1Dot0:
		if d.ServerHasCapability(v1Dot0Cap) {
			d.SelectedVersion = V1Dot0
		} else {
			return fmt.Errorf(
				"%w: user requested netconf version 1.0,"+
					" but server does not support this capability",
				util.ErrNetconfError,
			)
		}
	case V1Dot1:
		if d.ServerHasCapability(v1Dot1Cap) {
			d.SelectedVersion = V1Dot1
		} else {
			return fmt.Errorf(
				"%w: user requested netconf version 1.1,"+
					" but server does not support this capability",
				util.ErrNetconfError)
		}
	}

	ncPatterns := getNetconfPatterns()

	switch d.SelectedVersion {
	case V1Dot0:
		d.Channel.PromptPattern = ncPatterns.v1Dot0Delim
	case V1Dot1:
		d.Channel.PromptPattern = ncPatterns.v1Dot1Delim
	}

	return nil
}

func (d *Driver) sendClientCapabilities() error {
	var caps []byte

	switch d.SelectedVersion {
	case V1Dot0:
		caps = []byte(v1Dot0Caps)
	case V1Dot1:
		caps = []byte(v1Dot1Caps)
	}

	err := d.Channel.WriteAndReturn(caps, false)
	if err != nil {
		return err
	}

	return nil
}
