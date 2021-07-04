package netconf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/logging"
)

func (c *Channel) sendInput(channelInput []byte, stripPrompt, eager bool) *channelResult {
	logging.LogDebug(
		c.BaseChannel.FormatLogMessage(
			"info",
			fmt.Sprintf(
				"\"sending channelInput: %s; stripPrompt: %t; eager: %v",
				channelInput,
				stripPrompt,
				eager,
			),
		),
	)

	var b []byte

	err := c.BaseChannel.Write(channelInput, false)
	if err != nil {
		return &channelResult{
			result: []byte(""),
			error:  nil,
		}
	}

	err = c.readUntilInput(channelInput)
	if err != nil {
		return &channelResult{
			result: []byte(""),
			error:  err,
		}
	}

	err = c.BaseChannel.SendReturn()
	if err != nil {
		return &channelResult{
			result: []byte(""),
			error:  err,
		}
	}

	if !eager {
		postInputBuf, readErr := c.readUntilPrompt(b, nil)

		if readErr != nil {
			return &channelResult{
				result: []byte(""),
				error:  readErr,
			}
		}

		b = append(b, postInputBuf...)
	}

	return &channelResult{
		result: c.BaseChannel.RestructureOutput(b, stripPrompt),
		error:  nil,
	}
}

// SendInputBytes same as SendInput but accepting bytes straight away (this is mostly used for
//  netconf).
func (c *Channel) SendInputBytes(
	channelInput []byte,
	stripPrompt, eager bool,
	timeoutOps time.Duration,
) ([]byte, error) {
	_c := make(chan *channelResult)

	go func() {
		r := c.sendInput(channelInput, stripPrompt, eager)
		_c <- r
		close(_c)
	}()

	timer := time.NewTimer(c.BaseChannel.DetermineOperationTimeout(timeoutOps))

	select {
	case r := <-_c:
		return r.result, r.error
	case <-timer.C:
		logging.LogError(
			c.BaseChannel.FormatLogMessage("error", "timed out sending input to device"),
		)

		return []byte{}, channel.ErrChannelTimeout
	}
}

// SendInputNetconf like channel's `SendInput` but specifically for netconf; operates in eager mode
// unlike "normal" send input operations.
func (c *Channel) SendInputNetconf(channelInput []byte) ([]byte, error) {
	b, err := c.SendInputBytes(channelInput, false, true, -1)
	if err != nil {
		return b, err
	}

	if bytes.Contains(b, channelInput) {
		b = bytes.Split(b, channelInput)[1]
	}

	b, err = c.readUntilPrompt(b, nil)
	if err != nil {
		return b, err
	}

	if c.serverEcho == nil {
		logging.LogDebug(c.BaseChannel.FormatLogMessage(
			"debug", "server echo is unset, determining if server echoes inputs now"),
		)

		if bytes.Contains(b, channelInput) {
			logging.LogDebug(c.BaseChannel.FormatLogMessage(
				"info", "server echoes inputs, setting serverEcho to 'true'"),
			)

			echo := true
			c.serverEcho = &echo

			b, err = c.readUntilPrompt([]byte{}, nil)
			if err != nil {
				return b, err
			}
		} else {
			logging.LogDebug(c.BaseChannel.FormatLogMessage(
				"info", "server does *not* echo inputs, setting serverEcho to 'false'"),
			)

			echo := false
			c.serverEcho = &echo
		}
	}

	if c.SelectedNetconfVersion == Version11 {
		returnErr := c.BaseChannel.SendReturn()
		if returnErr != nil {
			return b, returnErr
		}
	}

	return b, nil
}
