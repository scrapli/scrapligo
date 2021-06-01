package channel

import (
	"fmt"
	"time"

	"github.com/scrapli/scrapligo/logging"
)

// SendInteractiveEvent struct used to represent each iteration of a `SendInteractive` operation.
type SendInteractiveEvent struct {
	ChannelInput    string
	ChannelResponse string
	HideInput       bool
}

func (c *Channel) sendInteractive(events []*SendInteractiveEvent) *channelResult {
	var b []byte

	for _, event := range events {
		channelInput := []byte(event.ChannelInput)
		channelResponse := event.ChannelResponse
		hideInput := event.HideInput

		logging.LogDebug(
			c.FormatLogMessage(
				"info",
				fmt.Sprintf(
					"\"sending interactive input: %s; expecting: %s; hidden input: %v",
					channelInput,
					channelResponse,
					hideInput,
				),
			),
		)

		err := c.Write(channelInput, hideInput)
		if err != nil {
			return &channelResult{
				result: []byte(""),
				error:  err,
			}
		}

		if channelResponse == "" || hideInput {
			returnErr := c.SendReturn()
			if returnErr != nil {
				return &channelResult{
					result: []byte(""),
					error:  returnErr,
				}
			}
		} else {
			newBuf, readErr := c.readUntilInput(channelInput)
			if readErr != nil {
				return &channelResult{
					result: []byte(""),
					error:  readErr,
				}
			}
			b = append(b, newBuf...)
			returnErr := c.SendReturn()
			if returnErr != nil {
				return &channelResult{
					result: []byte(""),
					error:  returnErr,
				}
			}
		}

		postInputBuf, err := c.readUntilPrompt(&channelResponse)
		if err != nil {
			return &channelResult{
				result: []byte(""),
				error:  nil,
			}
		}

		b = append(b, postInputBuf...)
	}

	return &channelResult{
		result: c.RestructureOutput(b, false),
		error:  nil,
	}
}

// SendInteractive send "interactive" inputs to a device. Accepts a slice of `SendInteractiveEvent`
// which is basically a struct defining the input and what the expected output of that command is.
// Used for dealing with "prompting" from a target device.
func (c *Channel) SendInteractive(
	events []*SendInteractiveEvent,
	timeoutOps time.Duration,
) ([]byte, error) {
	_c := make(chan *channelResult, 1)

	go func() {
		r := c.sendInteractive(events)
		_c <- r
		close(_c)
	}()

	timer := time.NewTimer(c.DetermineOperationTimeout(timeoutOps))

	select {
	case r := <-_c:
		return r.result, r.error
	case <-timer.C:
		logging.LogError(
			c.FormatLogMessage("error", "timed out sending interactive input to device"),
		)

		return []byte{}, ErrChannelTimeout
	}
}
