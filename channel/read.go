package channel

import (
	"bytes"
	"context"
	"errors"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/scrapli/scrapligo/util"
)

const inputSearchDepthMultiplier = 2

func getProcessReadBufSearchDepth(promptSearchDepth, inputLen int) int {
	finalSearchDepth := promptSearchDepth

	possibleSearchDepth := inputSearchDepthMultiplier * inputLen

	if possibleSearchDepth > finalSearchDepth {
		finalSearchDepth = possibleSearchDepth
	}

	return finalSearchDepth
}

func processReadBuf(rb []byte, searchDepth int) []byte {
	if len(rb) <= searchDepth {
		return rb
	}

	prb := rb[len(rb)-searchDepth:]

	partitionIdx := bytes.Index(prb, []byte("\n"))

	if partitionIdx > 0 {
		prb = prb[partitionIdx:]
	}

	return prb
}

func (c *Channel) read() {
	defer close(c.readDone)

	for {
		select {
		case <-c.done:
			return
		default:
		}

		b, err := c.t.Read()
		if err != nil {
			select {
			case <-c.done:
				// if the channel is shutting down, exit immediately -- there is no point
				// in reporting transport errors that are a side-effect of the close
				c.l.Debugf("discarding transport read error during shutdown: %s", err)
				return
			default:
			}

			if errors.Is(err, io.EOF) {
				// the underlying transport was closed so just return, we *probably* will have
				// already bailed out by reading from the (maybe/probably) closed done channel, but
				// if we hit EOF we know we are done anyway
				return
			}

			if strings.Contains(err.Error(), "input/output error") {
				// on Linux (including containerized environments) closing a PTY master fd
				// returns EIO ("input/output error") rather than io.EOF -- this is a known
				// kernel behavior documented in the system transport tests. treat it the same
				// as EOF: the underlying transport is gone and we are done
				return
			}

			// we got a transport error, put it into the error channel for processing during
			// the next read activity, log it, sleep and then try again...
			c.l.Criticalf(
				"encountered error reading from transport during channel read loop. error: %s", err,
			)

			select {
			case c.Errs <- err:
			case <-c.done:
				c.l.Debugf("discarding transport error during shutdown: %s", err)

				return
			}

			time.Sleep(c.ReadDelay)

			continue
		}

		if len(b) == 0 {
			// nothing to process... no reason to enqueue empty bytes, sleep and then continue...
			time.Sleep(c.ReadDelay)

			continue
		}

		// not 100% this is required, but has existed in scrapli/scrapligo for a long time and am
		// afraid to remove it!
		b = bytes.ReplaceAll(b, []byte("\r"), []byte(""))

		if bytes.Contains(b, []byte("\x1b")) {
			b = util.StripANSI(b)
		}

		c.Q.Enqueue(b)

		if c.ChannelLog != nil {
			_, err = c.ChannelLog.Write(b)
			if err != nil {
				c.l.Criticalf("error writing to channel log, ignoring. error: %s", err)
			}
		}

		time.Sleep(c.ReadDelay)
	}
}

// Read reads and returns the first available bytes from the channel Q object. If there are any
// errors on the Errs channel (these would come from the underlying transport), the error is
// returned with nil for the byte slice.
func (c *Channel) Read() ([]byte, error) {
	select {
	case <-c.readDone:
		return nil, util.ErrConnectionError
	case err, ok := <-c.Errs:
		if !ok {
			return nil, util.ErrConnectionError
		}

		return nil, err
	default:
	}

	b := c.Q.Dequeue()

	if b == nil {
		return nil, nil
	}

	c.l.Debugf("channel read %#v", string(b))

	return b, nil
}

// ReadAll reads and returns *all* available bytes form the channel Q object. If there are any
// errors on the Errs channel  (these would come from the underlying transport), the error is
// returned with nil for the byte slice. Be careful using this as it is possible to dequeue "too
// much" from the channel causing us to not be able to "find" the prompt or inputs during normal
// operations. In general, this should probably only be used when connecting to consoles/files.
func (c *Channel) ReadAll() ([]byte, error) {
	select {
	case <-c.readDone:
		return nil, util.ErrConnectionError
	case err, ok := <-c.Errs:
		if !ok {
			return nil, util.ErrConnectionError
		}

		return nil, err
	default:
	}

	b := c.Q.DequeueAll()

	if b == nil {
		return nil, nil
	}

	c.l.Debugf("channel read %#v", string(b))

	return b, nil
}

// ReadUntilFuzzy reads until a fuzzy match of the input is found.
func (c *Channel) ReadUntilFuzzy(ctx context.Context, b []byte) ([]byte, error) {
	if len(b) == 0 {
		return nil, nil
	}

	var rb []byte

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		nb, err := c.Read()
		if err != nil {
			return nil, err
		}

		if nb == nil {
			time.Sleep(c.ReadDelay)

			continue
		}

		rb = append(rb, nb...)

		if util.BytesRoughlyContains(
			b,
			processReadBuf(rb, getProcessReadBufSearchDepth(c.PromptSearchDepth, len(b))),
		) {
			return rb, nil
		}
	}
}

// ReadUntilExplicit reads bytes out of the channel Q object until the bytes b are seen in the
// output. Once the bytes are seen all read bytes are returned.
func (c *Channel) ReadUntilExplicit(ctx context.Context, b []byte) ([]byte, error) {
	var rb []byte

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		nb, err := c.Read()
		if err != nil {
			return nil, err
		}

		if nb == nil {
			time.Sleep(c.ReadDelay)

			continue
		}

		rb = append(rb, nb...)

		if bytes.Contains(
			processReadBuf(rb, getProcessReadBufSearchDepth(c.PromptSearchDepth, len(b))),
			b,
		) {
			return rb, nil
		}
	}
}

// ReadUntilPrompt reads bytes out of the channel Q object until the channel PromptPattern regex
// pattern is seen in the output. Once that pattern is seen, all read bytes are returned.
func (c *Channel) ReadUntilPrompt(ctx context.Context) ([]byte, error) {
	var rb []byte

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		nb, err := c.Read()
		if err != nil {
			return nil, err
		}

		if nb == nil {
			time.Sleep(c.ReadDelay)

			continue
		}

		rb = append(rb, nb...)

		if c.PromptPattern.Match(processReadBuf(rb, c.PromptSearchDepth)) {
			return rb, nil
		}
	}
}

// ReadUntilAnyPrompt reads bytes out of the channel Q object until any of the prompts in the
// "prompts" argument are seen in the output. Once any pattern is seen, all read bytes are returned.
func (c *Channel) ReadUntilAnyPrompt(
	ctx context.Context,
	prompts []*regexp.Regexp,
) ([]byte, error) {
	var rb []byte

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		nb, err := c.Read()
		if err != nil {
			return nil, err
		}

		if nb == nil {
			time.Sleep(c.ReadDelay)

			continue
		}

		rb = append(rb, nb...)

		prb := processReadBuf(rb, c.PromptSearchDepth)

		for _, p := range prompts {
			if p.Match(prb) {
				return rb, nil
			}
		}
	}
}
