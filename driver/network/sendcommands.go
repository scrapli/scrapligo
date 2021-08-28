package network

import (
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/util"
)

// SendCommands basically the same as the base driver flavor, but acquires the
// `DefaultDesiredPriv` prior to sending the command.
func (d *Driver) SendCommands(
	c []string,
	o ...base.SendOption,
) (*base.MultiResponse, error) {
	finalOpts := d.ParseSendOptions(o)

	if d.CurrentPriv != d.DefaultDesiredPriv {
		err := d.AcquirePriv(d.DefaultDesiredPriv)
		if err != nil {
			return nil, err
		}
	}

	return d.Driver.FullSendCommands(
		c,
		finalOpts.FailedWhenContains,
		finalOpts.StripPrompt,
		finalOpts.StopOnFailed,
		finalOpts.Eager,
		finalOpts.TimeoutOps,
	)
}

// SendCommandsFromFile basically the same as the base driver flavor, but acquires the
// `DefaultDesiredPriv` prior to sending the command.
func (d *Driver) SendCommandsFromFile(
	f string,
	o ...base.SendOption,
) (*base.MultiResponse, error) {
	c, err := util.LoadFileLines(f)
	if err != nil {
		return nil, err
	}

	return d.SendCommands(c, o...)
}
