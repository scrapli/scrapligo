package channel

import (
	"time"

	"github.com/scrapli/scrapligo/logging"
)

func (c *Channel) getPrompt() *channelResult {
	err := c.SendReturn()
	if err != nil {
		return &channelResult{
			error: err,
		}
	}

	var b []byte

	for {
		chunk, readErr := c.Read()
		b = append(b, chunk...)

		if readErr != nil {
			return &channelResult{
				error: readErr,
			}
		}

		channelMatch := c.CommsPromptPattern.Match(b)
		if channelMatch {
			return &channelResult{
				result: b,
				error:  nil,
			}
		}
	}
}

// GetPrompt fetch the current prompt.
func (c *Channel) GetPrompt() (string, error) {
	_c := make(chan *channelResult)

	go func() {
		r := c.getPrompt()
		_c <- r
		close(_c)
	}()

	timer := time.NewTimer(c.DetermineOperationTimeout(c.TimeoutOps))

	select {
	case r := <-_c:
		if r.error != nil {
			return "", r.error
		}

		return string(r.result), nil
	case <-timer.C:
		logging.LogError(c.FormatLogMessage("error", "timed out sending getting prompt"))

		return "", ErrChannelTimeout
	}
}
