package base

import (
	"bytes"
	"fmt"
	"regexp"
	"time"
)

type ReadCallback struct {
	Callback           func(*Driver, string) error
	Contains           string
	containsBytes      []byte
	ContainsRe         string
	containsReCompiled *regexp.Regexp
	CaseInsensitive    bool
	MultiLine          bool
	ResetOutput        bool
	Complete           bool
	Name               string
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

		r.containsReCompiled = regexp.MustCompile(fmt.Sprintf(`%s%s`, flags, r.contains()))
	}

	return r.containsReCompiled
}

func (d *Driver) ReadWithCallbacks( //nolint:gocognit
	callbacks []*ReadCallback,
	input string,
	output []byte,
	sleep time.Duration,
) error {
	if input != "" {
		err := d.Channel.WriteAndReturn([]byte(input), false)
		if err != nil {
			return err
		}

		return d.ReadWithCallbacks(callbacks, "", nil, sleep)
	}

	for {
		newOutput, err := d.Channel.Read()
		if err != nil {
			return err
		}

		output = append(output, newOutput...)

		for _, callback := range callbacks {
			o := output
			if callback.CaseInsensitive {
				o = bytes.ToLower(output)
			}

			if (callback.Contains != "" && bytes.Contains(o, callback.contains())) ||
				(callback.ContainsRe != "" && callback.containsRe().Match(o)) {
				err = callback.Callback(d, string(output))
				if err != nil {
					return err
				}

				if callback.Complete {
					return nil
				}

				if callback.ResetOutput {
					output = []byte{}
				}

				return d.ReadWithCallbacks(callbacks, "", output, sleep)
			}
		}

		time.Sleep(sleep)
	}
}
