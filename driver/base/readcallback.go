package base

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var ErrCallbackAlreadyTriggered = errors.New("callback set to 'OnlyOnce', but already triggered")

type ReadCallback struct {
	Callback           func(*Driver, string) error
	Contains           string
	containsBytes      []byte
	ContainsRe         string
	containsReCompiled *regexp.Regexp
	CaseInsensitive    bool
	MultiLine          bool
	// ResetOutput bool indicating if the output should be reset or not after callback execution.
	ResetOutput bool
	// OnlyOnce bool indicating if this callback should be executed only one time.
	OnlyOnce  bool
	triggered bool
	Complete  bool
	Name      string
}

func (r *ReadCallback) contains() []byte {
	if len(r.containsBytes) == 0 {
		r.containsBytes = []byte(r.Contains)
	}

	return r.containsBytes
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

func (d *Driver) executeCallback(
	i int,
	callbacks []*ReadCallback,
	output []byte,
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

	return d.readWithCallbacks(callbacks, output, readDelay)
}

func (d *Driver) readWithCallbacks(
	callbacks []*ReadCallback,
	output []byte,
	readDelay time.Duration,
) error {
	for {
		newOutput, err := d.Channel.Read()
		if err != nil {
			return err
		}

		output = append(output, newOutput...)

		for i, callback := range callbacks {
			o := output
			if callback.CaseInsensitive {
				o = bytes.ToLower(output)
			}

			if (callback.Contains != "" && bytes.Contains(o, callback.contains())) ||
				(callback.ContainsRe != "" && callback.containsRe().Match(o)) {
				return d.executeCallback(i, callbacks, output, readDelay)
			}
		}

		time.Sleep(readDelay)
	}
}

func (d *Driver) ReadWithCallbacks(
	callbacks []*ReadCallback,
	input string,
	readDelay time.Duration,
) error {
	if input != "" {
		err := d.Channel.WriteAndReturn([]byte(input), false)
		if err != nil {
			return err
		}
	}

	return d.readWithCallbacks(callbacks, []byte{}, readDelay)
}
