package netconf

import (
	"bytes"
	"regexp"
	"time"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/channel"
)

// Channel the Netconf Channel that extends the base SSH channel.
type Channel struct {
	BaseChannel             *channel.Channel
	PreferredNetconfVersion string
	SelectedNetconfVersion  string
	serverCapabilities      []string
	serverEcho              *bool
}

type channelResult struct {
	result []byte
	error  error
}

// OpenNetconf open a netconf channel to the device, handles capabilities exchange.
func (c *Channel) OpenNetconf(authenticationBuf []byte) error {
	if !bytes.Contains(authenticationBuf, []byte("]]>]]>")) {
		var _c = make(chan *channelResult, 1)

		go func() {
			r := c.getServerCapabilities(authenticationBuf)
			_c <- r
			close(_c)
		}()

		timer := time.NewTimer(*c.BaseChannel.TimeoutOps)

		select {
		case r := <-_c:
			authenticationBuf = r.result
		case <-timer.C:
			logging.LogError(
				c.BaseChannel.FormatLogMessage(
					"error",
					"timed out attempting to read server capabilities",
				),
			)

			return channel.ErrAuthTimeout
		}
	}

	err := c.parseServerCapabilities(authenticationBuf)
	if err != nil {
		return err
	}

	err = c.processCapabilitiesExchange()
	if err != nil {
		return err
	}

	err = c.sendClientCapabilities()
	if err != nil {
		return err
	}

	return nil
}

func (c *Channel) readUntilInput(channelInput []byte) error {
	var b []byte

	if !*c.serverEcho || len(channelInput) == 0 {
		return nil
	}

	for {
		chunk, err := c.BaseChannel.Read()

		b = append(b, chunk...)

		if err != nil {
			return err
		}

		if bytes.Contains(b, channelInput) || bytes.Contains(b, []byte("rpc>")) {
			return err
		}
	}
}

func (c *Channel) readUntilPrompt(b []byte, prompt *string) ([]byte, error) {
	matchPattern := c.BaseChannel.CommsPromptPattern
	if prompt != nil {
		matchPattern = regexp.MustCompile(*prompt)
	}

	for {
		chunk, err := c.BaseChannel.Read()
		b = append(b, chunk...)

		if err != nil {
			return b, err
		}

		channelMatch := matchPattern.Match(b)
		if channelMatch {
			logging.LogDebug(c.BaseChannel.FormatLogMessage("debug", "found prompt match"))

			return b, err
		}
	}
}

func (c *Channel) checkEcho() error {
	var _c = make(chan error, 1)

	echoTimeout := 1 * time.Second
	if *c.BaseChannel.TimeoutOps > 0*time.Second {
		echoTimeout = *c.BaseChannel.TimeoutOps / 20
	}

	go func() {
		// try to read a single byte off the transport
		_, err := c.BaseChannel.Transport.ReadN(1)

		_c <- err
		close(_c)
	}()

	timer := time.NewTimer(echoTimeout)

	select {
	case err := <-_c:
		if err != nil {
			return err
		}

		logging.LogDebug(c.BaseChannel.FormatLogMessage(
			"info", "server echoes inputs, setting serverEcho to 'true'"),
		)

		echo := true

		c.serverEcho = &echo
	case <-timer.C:
		logging.LogDebug(c.BaseChannel.FormatLogMessage(
			"info", "server does *not* echo inputs, setting serverEcho to 'false'"),
		)

		echo := false

		c.serverEcho = &echo
	}

	return nil
}
