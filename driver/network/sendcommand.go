package network

import (
	"github.com/scrapli/scrapligo/driver/base"
)

// SendCommand basically the same as the base driver flavor, but acquires the
// `DefaultDesiredPriv` prior to sending the command.
func (d *Driver) SendCommand(c string, o ...base.SendOption) (*base.Response, error) {
	finalOpts := d.ParseSendOptions(o)

	if d.CurrentPriv != d.DefaultDesiredPriv {
		err := d.AcquirePriv(d.DefaultDesiredPriv)
		if err != nil {
			return nil, err
		}
	}

	return d.Driver.FullSendCommand(
		c,
		finalOpts.FailedWhenContains,
		finalOpts.StripPrompt,
		finalOpts.Eager,
		finalOpts.TimeoutOps,
	)
}
