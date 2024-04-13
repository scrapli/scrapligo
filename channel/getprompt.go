package channel

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapli/scrapligo/util"
)

// GetPrompt returns a byte slice containing the current "prompt" of the connected ssh/telnet
// server.
func (c *Channel) GetPrompt() ([]byte, error) {
	c.l.Debug("channel GetPrompt requested")

	cr := make(chan *result)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func() {
		err := c.WriteReturn()
		if err != nil {
			cr <- &result{b: nil, err: err}

			return
		}

		var b []byte

		b, err = c.ReadUntilPrompt(ctx)

		// we already know the pattern is in the buf, we just want ot re to yoink it out without
		// any newlines or extra stuff we read (which shouldn't happen outside the initial
		// connection but...)
		cr <- &result{b: c.PromptPattern.Find(b), err: err}
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
