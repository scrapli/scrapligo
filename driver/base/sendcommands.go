package base

import (
	"time"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/logging"
)

// FullSendCommands same as `SendCommands` but requiring explicit options.
func (d *Driver) FullSendCommands(
	c []string,
	failedWhenContains []string,
	stripPrompt, stopOnFailed, eager bool,
	timeoutOps time.Duration,
) (*MultiResponse, error) {
	mr := NewMultiResponse(d.Host)

	for _, command := range c[:len(c)-1] {
		r, err := d.FullSendCommand(
			command,
			failedWhenContains,
			stripPrompt,
			eager,
			timeoutOps,
		)

		mr.AppendResponse(r)

		if err != nil {
			return mr, err
		}

		if stopOnFailed && r.Failed != nil {
			logging.LogDebug(
				d.FormatLogMessage(
					"info",
					"encountered failed command, and stop on failed is true, discontinuing send"+
						" commands operation",
				),
			)

			return mr, err
		}
	}

	r, err := d.FullSendCommand(
		c[len(c)-1],
		failedWhenContains,
		stripPrompt,
		eager,
		timeoutOps,
	)
	mr.AppendResponse(r)

	return mr, err
}

// SendCommands send commands to a device, accepts a string command and variadic of `SendOption`s.
func (d *Driver) SendCommands(
	c []string,
	o ...SendOption,
) (*MultiResponse, error) {
	finalOpts := d.ParseSendOptions(o)

	return d.FullSendCommands(
		c,
		finalOpts.FailedWhenContains,
		finalOpts.StripPrompt,
		finalOpts.StopOnFailed,
		finalOpts.Eager,
		finalOpts.TimeoutOps,
	)
}

// SendCommandsFromFile send commands from a file to a device, accepts a string command and variadic
// of `SendOption`s.
func (d *Driver) SendCommandsFromFile(
	f string,
	o ...SendOption,
) (*MultiResponse, error) {
	finalOpts := d.ParseSendOptions(o)

	c, err := util.LoadFileLines(f)
	if err != nil {
		return nil, err
	}

	return d.FullSendCommands(
		c,
		finalOpts.FailedWhenContains,
		finalOpts.StripPrompt,
		finalOpts.StopOnFailed,
		finalOpts.Eager,
		finalOpts.TimeoutOps,
	)
}
