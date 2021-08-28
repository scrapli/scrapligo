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
	ForceSelfClosingTag     bool
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
		var _c = make(chan *channelResult)

		go func() {
			r := c.getServerCapabilities(authenticationBuf)
			_c <- r
			close(_c)
		}()

		timer := time.NewTimer(c.BaseChannel.TimeoutOps)

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

	if c.serverEcho == nil {
		// echo check hasnt been figured out yet
		return nil
	}

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
