package channel

import (
	"fmt"
	"time"

	"github.com/scrapli/scrapligo/util"
)

// GetPrompt returns a byte slice containing the current "prompt" of the connected ssh/telnet
// server.
func (c *Channel) GetPrompt() ([]byte, error) {
	c.l.Debug("channel GetPrompt requested")

	cr := make(chan *result)

	go func() {
		err := c.WriteReturn()
		if err != nil {
			cr <- &result{b: nil, err: err}

			return
		}

		var b []byte

		b, err = c.ReadUntilPrompt()

		cr <- &result{b: b, err: err}
	}()

	timer := time.NewTimer(c.TimeoutOps)

	select {
	case r := <-cr:
		if r.err != nil {
			return nil, r.err
		}

		return r.b, nil
	case <-timer.C:
		c.l.Critical("channel timeout fetching prompt")

		return nil, fmt.Errorf("%w: channel timeout fetching prompt", util.ErrTimeoutError)
	}
}
