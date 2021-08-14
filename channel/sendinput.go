package channel

import (
	"fmt"
	"time"

	"github.com/scrapli/scrapligo/logging"
)

func (c *Channel) sendInput(channelInput []byte, stripPrompt, eager bool) *channelResult {
	logging.LogDebug(
		c.FormatLogMessage(
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

	err := c.Write(channelInput, false)
	if err != nil {
		return &channelResult{
			result: []byte(""),
			error:  nil,
		}
	}

	_, err = c.readUntilInput(channelInput)
	if err != nil {
		return &channelResult{
			result: []byte(""),
			error:  err,
		}
	}

	err = c.SendReturn()
	if err != nil {
		return &channelResult{
			result: []byte(""),
			error:  err,
		}
	}

	if !eager {
		postInputBuf, readErr := c.readUntilPrompt()

		if readErr != nil {
			return &channelResult{
				result: []byte(""),
				error:  readErr,
			}
		}

		b = append(b, postInputBuf...)
	}

	return &channelResult{
		result: c.RestructureOutput(b, stripPrompt),
		error:  nil,
	}
}

// SendInput send input to the device, optionally strips prompt out of the returned output. Eager
// flag should generally not be used unless you know what you are doing! The `timeoutOps` argument
//  modifies the base timeout argument just for the duration of this send operation.
func (c *Channel) SendInput(
	channelInput string,
	stripPrompt, eager bool,
	timeoutOps time.Duration,
) ([]byte, error) {
	_c := make(chan *channelResult)

	go func() {
		r := c.sendInput([]byte(channelInput), stripPrompt, eager)
		_c <- r
		close(_c)
	}()

	timer := time.NewTimer(c.DetermineOperationTimeout(timeoutOps))

	select {
	case r := <-_c:
		return r.result, r.error
	case <-timer.C:
		logging.LogError(c.FormatLogMessage("error", "timed out sending input to device"))
		return []byte{}, ErrChannelTimeout
	}
}
