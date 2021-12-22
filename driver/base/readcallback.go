package base

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var ErrMustSetContains = errors.New("must set Contains or ContainsRe")
var ErrCallbackAlreadyTriggered = errors.New("callback set to 'OnlyOnce', but already triggered")
var ErrCallbackTimeout = errors.New("callback timeout")

type ReadCallbackOption func(callback *ReadCallback) error

func WithCallbackContains(contains string) ReadCallbackOption {
	return func(r *ReadCallback) error {
		r.Contains = contains

		return nil
	}
}

func WithCallbackNotContains(notContains string) ReadCallbackOption {
	return func(r *ReadCallback) error {
		r.NotContains = notContains

		return nil
	}
}

func WithCallbackContainsRe(contains string) ReadCallbackOption {
	return func(r *ReadCallback) error {
		r.ContainsRe = contains

		return nil
	}
}

func WithCallbackCaseInsensitive(i bool) ReadCallbackOption {
	return func(r *ReadCallback) error {
		r.CaseInsensitive = i

		return nil
	}
}

func WithCallbackMultiline(m bool) ReadCallbackOption {
	return func(r *ReadCallback) error {
		r.MultiLine = m

		return nil
	}
}

func WithCallbackResetOutput(reset bool) ReadCallbackOption {
	return func(r *ReadCallback) error {
		r.ResetOutput = reset

		return nil
	}
}

func WithCallbackOnlyOnce(o bool) ReadCallbackOption {
	return func(r *ReadCallback) error {
		r.OnlyOnce = o

		return nil
	}
}

func WithCallbackNextTimeout(t time.Duration) ReadCallbackOption {
	return func(r *ReadCallback) error {
		r.NextTimeout = t

		return nil
	}
}

func WithCallbackNextReadDelay(t time.Duration) ReadCallbackOption {
	return func(r *ReadCallback) error {
		r.NextReadDelay = t

		return nil
	}
}

func WithCallbackComplete(complete bool) ReadCallbackOption {
	return func(r *ReadCallback) error {
		r.Complete = complete

		return nil
	}
}

func WithCallbackName(name string) ReadCallbackOption {
	return func(r *ReadCallback) error {
		r.Name = name

		return nil
	}
}

func NewReadCallback(
	callback func(*Driver, string) error,
	options ...ReadCallbackOption,
) (*ReadCallback, error) {
	rc := &ReadCallback{
		Callback:           callback,
		Contains:           "",
		containsBytes:      nil,
		ContainsRe:         "",
		containsReCompiled: nil,
		CaseInsensitive:    true,
		MultiLine:          true,
		ResetOutput:        true,
		OnlyOnce:           false,
		NextTimeout:        0,
		NextReadDelay:      0,
		triggered:          false,
		Complete:           false,
		Name:               "",
	}

	for _, option := range options {
		err := option(rc)
		if err != nil {
			return nil, err
		}
	}

	if rc.Contains == "" && rc.ContainsRe == "" {
		return nil, ErrMustSetContains
	}

	return rc, nil
}

type ReadCallback struct {
	Callback           func(*Driver, string) error
	Contains           string
	containsBytes      []byte
	NotContains        string
	notContainsBytes   []byte
	ContainsRe         string
	containsReCompiled *regexp.Regexp
	CaseInsensitive    bool
	MultiLine          bool
	// ResetOutput bool indicating if the output should be reset or not after callback execution.
	ResetOutput bool
	// OnlyOnce bool indicating if this callback should be executed only one time.
	OnlyOnce bool
	// NextTimout timeout value to use for the subsequent read loop - ignored if Complete is true.
	NextTimeout time.Duration
	// NextReadDelay is time to use for sleeps between reads for hte subsequent read loop.
	NextReadDelay time.Duration
	triggered     bool
	Complete      bool
	Name          string
}

func (r *ReadCallback) contains() []byte {
	if len(r.containsBytes) == 0 {
		r.containsBytes = []byte(r.Contains)

		if r.CaseInsensitive {
			r.containsBytes = bytes.ToLower(r.containsBytes)
		}
	}

	return r.containsBytes
}

func (r *ReadCallback) notContains() []byte {
	if len(r.notContainsBytes) == 0 {
		r.notContainsBytes = []byte(r.NotContains)

		if r.CaseInsensitive {
			r.notContainsBytes = bytes.ToLower(r.notContainsBytes)
		}
	}

	return r.notContainsBytes
}

func (r *ReadCallback) containsRe() *regexp.Regexp {
	if r.containsReCompiled == nil {
		flags := ""

		if r.CaseInsensitive && r.MultiLine {
			flags = "(?im)"
		} else if r.CaseInsensitive {
			flags = "(?i)"
		} else if r.MultiLine {
			flags = "(?m)"
		}

		r.containsReCompiled = regexp.MustCompile(fmt.Sprintf(`%s%s`, flags, r.ContainsRe))
	}

	return r.containsReCompiled
}

func (r *ReadCallback) check(o []byte) bool {
	if r.CaseInsensitive {
		o = bytes.ToLower(o)
	}

	if (r.Contains != "" && bytes.Contains(o, r.contains())) &&
		!(r.NotContains != "" && !bytes.Contains(o, r.notContains())) {
		return true
	}

	if (r.ContainsRe != "" && r.containsRe().Match(o)) &&
		!(r.NotContains != "" && !bytes.Contains(o, r.notContains())) {
		return true
	}

	return false
}

type readCallbackResult struct {
	i         int
	callbacks []*ReadCallback
	output    []byte
	err       error
}

func (d *Driver) executeCallback(
	i int,
	callbacks []*ReadCallback,
	output []byte,
	timeout,
	readDelay time.Duration) error {
	callback := callbacks[i]

	if callback.OnlyOnce {
		if callback.triggered {
			return ErrCallbackAlreadyTriggered
		}

		callback.triggered = true
	}

	err := callback.Callback(d, string(output))
	if err != nil {
		return err
	}

	if callback.Complete {
		return nil
	}

	if callback.ResetOutput {
		output = []byte{}
	}

	nextTimeout := timeout
	if callback.NextTimeout != 0 {
		nextTimeout = callback.NextTimeout
	}

	nextReadDelay := readDelay
	if callback.NextReadDelay != 0 {
		nextReadDelay = callback.NextReadDelay
	}

	return d.readWithCallbacks(callbacks, output, nextTimeout, nextReadDelay)
}

func (d *Driver) readWithCallbacks(
	callbacks []*ReadCallback,
	output []byte,
	timeout,
	readDelay time.Duration,
) error {
	c := make(chan *readCallbackResult)

	go func() {
		defer close(c)

		for {
			newOutput, err := d.Channel.Read()
			if err != nil {
				c <- &readCallbackResult{
					err: err,
				}

				return
			}

			output = append(output, newOutput...)

			for i, callback := range callbacks {
				if callback.check(output) {
					c <- &readCallbackResult{
						i:         i,
						callbacks: callbacks,
						output:    output,
						err:       nil,
					}

					return
				}
			}

			time.Sleep(readDelay)
		}
	}()

	timer := time.NewTimer(timeout)

	select {
	case r := <-c:
		if r.err != nil {
			return r.err
		}

		return d.executeCallback(r.i, r.callbacks, r.output, timeout, readDelay)
	case <-timer.C:
		return ErrCallbackTimeout
	}
}

func (d *Driver) ReadWithCallbacks(
	callbacks []*ReadCallback,
	input string,
	timeout,
	readDelay time.Duration,
) error {
	if input != "" {
		err := d.Channel.WriteAndReturn([]byte(input), false)
		if err != nil {
			return err
		}
	}

	origTransportTimeout := d.Transport.BaseTransportArgs.TimeoutTransport
	d.Transport.BaseTransportArgs.TimeoutTransport = 0

	r := d.readWithCallbacks(callbacks, []byte{}, timeout, readDelay)

	d.Transport.BaseTransportArgs.TimeoutTransport = origTransportTimeout

	return r
}
