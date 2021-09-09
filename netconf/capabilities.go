package netconf

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/scrapli/scrapligo/logging"
)

func (c *Channel) getServerCapabilities(authenticationBuf []byte) *channelResult {
	for {
		b, err := c.BaseChannel.Read()
		if err != nil {
			return &channelResult{
				result: authenticationBuf,
				error:  err,
			}
		}

		authenticationBuf = append(authenticationBuf, b...)

		if bytes.Contains(authenticationBuf, []byte("]]>]]>")) {
			return &channelResult{
				result: authenticationBuf,
				error:  nil,
			}
		}
	}
}

func (c *Channel) parseServerCapabilities(authenticationBuf []byte) error {
	// match w/ or w/out hello namespace
	serverHelloPattern := regexp.MustCompile(`(?is)(<(\w+:)?hello.*</(\w+:)?hello>)`)
	serverHelloMatch := serverHelloPattern.Match(authenticationBuf)

	if !serverHelloMatch {
		logging.LogError(
			c.BaseChannel.FormatLogMessage("error", "could not parse server capabilities"),
		)

		return ErrCapabilitiesExchangeFailed
	}

	// rather than deal w/ xml like scrapli python does, just regex the caps out
	serverCapabilitiesPattern := regexp.MustCompile(
		`(?i)(?:<(?:\w+:)?capability>)(.*?)(?:</(?:\w+:)?capability>)`,
	)
	serverCapabilitiesMatches := serverCapabilitiesPattern.FindAllSubmatch(authenticationBuf, -1)

	serverCapabilities := make([]string, 1)
	for _, match := range serverCapabilitiesMatches {
		serverCapabilities = append(serverCapabilities, string(match[1]))
	}

	c.serverCapabilities = serverCapabilities

	return nil
}

func (c *Channel) serverCapabilitiesContains(requestedCapability string) bool {
	for _, serverCapability := range c.serverCapabilities {
		if serverCapability == requestedCapability {
			return true
		}
	}

	return false
}

func (c *Channel) processCapabilitiesExchange() error {
	if c.serverCapabilitiesContains(Version11Capability) {
		c.SelectedNetconfVersion = Version11
	} else if c.serverCapabilitiesContains(Version10Capability) {
		c.SelectedNetconfVersion = Version10
	} else {
		logging.LogError(c.BaseChannel.FormatLogMessage(
			"error", "did not receive netconf capabilities for version 1.0 or 1.1"),
		)
		return ErrCapabilitiesExchangeFailed
	}

	if c.PreferredNetconfVersion != "" {
		if c.PreferredNetconfVersion == Version10 &&
			c.serverCapabilitiesContains(Version10Capability) {
			c.SelectedNetconfVersion = Version10
		} else if c.PreferredNetconfVersion == Version11 &&
			c.serverCapabilitiesContains(Version11Capability) {
			c.SelectedNetconfVersion = Version11
		} else {
			logging.LogDebug(c.BaseChannel.FormatLogMessage(
				"info", "user provided preferred netconf version not available"),
			)
		}
	}

	if c.SelectedNetconfVersion == Version10 {
		c.BaseChannel.CommsPromptPattern = regexp.MustCompile(Version10DelimiterPattern)
	} else {
		c.BaseChannel.CommsPromptPattern = regexp.MustCompile(Version11DelimiterPattern)
	}

	return nil
}

func (c *Channel) sendClientCapabilities() error {
	clientCapabilities := Version11Capabilities
	if c.SelectedNetconfVersion == Version10 {
		clientCapabilities = Version10Capabilities
	}

	logging.LogDebug(c.BaseChannel.FormatLogMessage(
		"info", fmt.Sprintf("sending client capabilities: %s\n", clientCapabilities)),
	)

	err := c.BaseChannel.Write([]byte(clientCapabilities), false)
	if err != nil {
		return err
	}

	err = c.readUntilInput([]byte(clientCapabilities[1:]))
	if err != nil {
		return err
	}

	err = c.BaseChannel.SendReturn()
	if err != nil {
		return err
	}

	return nil
}
