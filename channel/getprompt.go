package channel

import (
	"bytes"
	"strings"
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

	b := make([]byte, 0, 100)

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
	_c := make(chan *channelResult, 1)

	go func() {
		r := c.getPrompt()
		_c <- r
		close(_c)
	}()

	timer := time.NewTimer(c.DetermineOperationTimeout(*c.TimeoutOps))

	select {
	case r := <-_c:
		if r.error != nil {
			return "", r.error
		}

		return strings.TrimSpace(string(bytes.Trim(r.result, "\x00"))), nil
	case <-timer.C:
		logging.LogError(c.FormatLogMessage("error", "timed out sending getting prompt"))

		return "", ErrChannelTimeout
	}
}
