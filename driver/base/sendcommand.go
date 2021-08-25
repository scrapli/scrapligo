package base

import (
	"time"

	"github.com/scrapli/scrapligo/response"
)

// FullSendCommand same as `SendCommand` but requiring explicit options.
func (d *Driver) FullSendCommand(
	c string,
	failedWhenContains []string,
	stripPrompt, eager bool,
	timeoutOps time.Duration,
) (*response.Response, error) {
	r := response.NewResponse(d.Host, d.Port, c, failedWhenContains)

	rawResult, err := d.Channel.SendInput(c, stripPrompt, eager, timeoutOps)

	r.Record(rawResult, string(rawResult))

	return r, err
}

// SendCommand send a command to a device, accepts a string command and variadic of `SendOption`s.
func (d *Driver) SendCommand(c string, o ...SendOption) (*response.Response, error) {
	finalOpts := d.ParseSendOptions(o)

	return d.FullSendCommand(
		c,
		finalOpts.FailedWhenContains,
		finalOpts.StripPrompt,
		finalOpts.Eager,
		finalOpts.TimeoutOps,
	)
}
