package generic

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/scrapli/scrapligo/channel"

	"github.com/scrapli/scrapligo/response"

	"github.com/scrapli/scrapligo/util"
)

// NewCallback returns a Callback object with provided options applied.
func NewCallback(
	callback func(*Driver, string) error,
	opts ...util.Option,
) (*Callback, error) {
	c := &Callback{
		Callback:      callback,
		Contains:      "",
		containsBytes: nil,
		ContainsRe:    nil,
		Insensitive:   true,
		ResetOutput:   true,
		Once:          false,
		NextTimeout:   0,
		triggered:     false,
		Complete:      false,
		Name:          "",
	}

	for _, option := range opts {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}

	if c.Contains == "" && c.ContainsRe == nil {
		return nil, fmt.Errorf("%w: must provide contains or contains regex", util.ErrBadOption)
	}

	return c, nil
}

// Callback represents not only a callback function, but what triggers that callback function. This
// object is only used in conjunction with the SendWithCallbacks method.
type Callback struct {
	Callback         func(*Driver, string) error
	Contains         string
	containsBytes    []byte
	NotContains      string
	notContainsBytes []byte
	ContainsRe       *regexp.Regexp
	Insensitive      bool
	// SearchDepth limits matching to newly read bytes plus this many preceding bytes.
	// Values <= 0 search the entire output since the last ResetOutput (the default).
	// Use this only for local markers whose required context fits within SearchDepth;
	// the full output is still passed to Callback and retained in the response.
	// NotContains requires the full segment and disables this optimization.
	SearchDepth int
	// ResetOutput bool indicating if the output should be reset or not after callback execution.
	ResetOutput bool
	// Once bool indicating if this callback should be executed only one time.
	Once bool
	// NextTimout timeout value to use for the subsequent read loop - ignored if Complete is true.
	NextTimeout time.Duration
	triggered   bool
	Complete    bool
	Name        string
}

func (c *Callback) contains() []byte {
	if len(c.containsBytes) == 0 {
		c.containsBytes = []byte(c.Contains)

		if c.Insensitive {
			c.containsBytes = bytes.ToLower(c.containsBytes)
		}
	}

	return c.containsBytes
}

func (c *Callback) notContains() []byte {
	if len(c.notContainsBytes) == 0 {
		c.notContainsBytes = []byte(c.NotContains)

		if c.Insensitive {
			c.notContainsBytes = bytes.ToLower(c.notContainsBytes)
		}
	}

	return c.notContainsBytes
}

func (c *Callback) check(b []byte, containsRe *regexp.Regexp) bool {
	if c.Insensitive {
		b = bytes.ToLower(b)
	}

	if (c.Contains != "" && bytes.Contains(b, c.contains())) &&
		!(c.NotContains != "" && !bytes.Contains(b, c.notContains())) {
		return true
	}

	if (containsRe != nil && containsRe.Match(b)) &&
		!(c.NotContains != "" && !bytes.Contains(b, c.notContains())) {
		return true
	}

	return false
}

// callbackSearch preserves the character before a bounded search window. Requiring
// the regex to consume that character keeps ^, \A and word boundaries from treating
// the cut in the buffer as the beginning of the original output.
type callbackSearch struct {
	callback   *Callback
	expression string
	prepared   bool
	pattern    *regexp.Regexp
}

func newCallbackSearch(cb *Callback) callbackSearch {
	return callbackSearch{callback: cb}
}

// boundedPattern is prepared lazily and refreshed if a callback changes its regex
// between stages. Cache failures too: wrapping a valid regex can exceed regexp's
// nesting limit, in which case the caller must keep full-segment matching.
func (s *callbackSearch) boundedPattern() *regexp.Regexp {
	if s.callback.ContainsRe == nil {
		return nil
	}

	expression := s.callback.ContainsRe.String()
	if !s.prepared || s.expression != expression {
		s.expression = expression
		s.prepared = true

		pattern, err := regexp.Compile(`(?s:.)(?:` + expression + `)`)
		if err != nil {
			s.pattern = nil
		} else {
			s.pattern = pattern
		}
	}

	return s.pattern
}

func (s *callbackSearch) check(b []byte, previousLen int) bool {
	cb := s.callback
	if cb.SearchDepth <= 0 || cb.NotContains != "" || previousLen <= cb.SearchDepth {
		return cb.check(b, cb.ContainsRe)
	}

	pattern := s.boundedPattern()
	if cb.ContainsRe != nil && pattern == nil {
		return cb.check(b, cb.ContainsRe)
	}

	start := previousLen - cb.SearchDepth
	// Do not split a UTF-8 character at the beginning of the search window.
	for shift := 0; shift < utf8.UTFMax-1 && start > 0 && !utf8.RuneStart(b[start]); shift++ {
		start--
	}

	if start == 0 {
		return cb.check(b, cb.ContainsRe)
	}

	_, size := utf8.DecodeLastRune(b[:start])

	return cb.check(b[start-size:], pattern)
}

func (d *Driver) executeCallback(cb *Callback, b []byte) error {
	if cb.Once {
		if cb.triggered {
			return fmt.Errorf(
				"%w: callback once set, and callback already triggered",
				util.ErrOperationError,
			)
		}

		cb.triggered = true
	}

	if cb.Callback != nil {
		// you might not want to set a callback on the "done" stage, so we skip executing if
		// callback is nil
		err := cb.Callback(d, string(b))
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Driver) handleCallbacks(
	callbacks []*Callback,
	b, fb []byte,
	timeout time.Duration,
) ([]byte, error) {
	searches := make([]callbackSearch, len(callbacks))
	for i, cb := range callbacks {
		searches[i] = newCallbackSearch(cb)
	}

	for {
		result, err := d.readCallback(searches, b, fb, timeout)
		if err != nil {
			return nil, err
		}

		cb := result.callback
		b, fb = result.b, result.fb

		if err := d.executeCallback(cb, b); err != nil {
			return nil, err
		}

		if cb.Complete {
			return fb, nil
		}

		if cb.ResetOutput {
			b = nil
		}

		if cb.NextTimeout != 0 {
			timeout = cb.NextTimeout
		}
	}
}

type callbackResult struct {
	callback *Callback
	b        []byte
	fb       []byte
}

// readCallback reads synchronously: Channel.Read is non-blocking, and no reader
// goroutine should remain on the channel after this operation times out.
func (d *Driver) readCallback(
	searches []callbackSearch,
	b, fb []byte,
	timeout time.Duration,
) (*callbackResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	readDelay := d.Channel.ReadDelay
	if readDelay <= 0 {
		readDelay = channel.DefaultReadDelayMicroSeconds * time.Microsecond
	}

	ticker := time.NewTicker(readDelay)
	defer ticker.Stop()

	// Check once when entering a callback stage, including a retained segment or
	// a pattern matching empty output. Otherwise only new data triggers matching.
	checkPending := true

	for {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("%w: timeout handling callbacks", util.ErrTimeoutError)
		}

		rb, err := d.Channel.Read()
		if err != nil {
			return nil, err
		}

		if len(rb) == 0 && !checkPending {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("%w: timeout handling callbacks", util.ErrTimeoutError)
			case <-ticker.C:
				continue
			}
		}

		previousLen := len(b)
		if checkPending {
			previousLen = 0
			checkPending = false
		}

		b = append(b, rb...)
		fb = append(fb, rb...)

		for i := range searches {
			search := &searches[i]
			matched := search.check(b, previousLen)

			if ctx.Err() != nil {
				return nil, fmt.Errorf("%w: timeout handling callbacks", util.ErrTimeoutError)
			}

			if !matched {
				continue
			}

			return &callbackResult{callback: search.callback, b: b, fb: fb}, nil
		}
	}
}

// SendWithCallbacks sends some input and responds to the output of that input based on the list of
// callbacks provided. This method can be looked at as a more advanced SendInteractive.
func (d *Driver) SendWithCallbacks(
	input string,
	callbacks []*Callback,
	timeout time.Duration,
	opts ...util.Option,
) (*response.Response, error) {
	d.Logger.Info("SendWithCallbacks requested")

	driverOpts, err := NewOperation(opts...)
	if err != nil {
		return nil, err
	}

	if len(driverOpts.FailedWhenContains) == 0 {
		driverOpts.FailedWhenContains = d.FailedWhenContains
	}

	r := response.NewResponse(
		input,
		d.Transport.GetHost(),
		d.Transport.GetPort(),
		driverOpts.FailedWhenContains,
	)

	if input != "" {
		err := d.Channel.WriteAndReturn([]byte(input), false)
		if err != nil {
			return nil, err
		}
	}

	b, err := d.handleCallbacks(callbacks, nil, nil, timeout)
	if err != nil {
		return nil, err
	}

	r.Record(b)

	return r, nil
}
