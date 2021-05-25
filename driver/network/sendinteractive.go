package network

import (
	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/driver/base"
)

// SendInteractive send interactive commands to a device, accepts a slice of `SendInteractiveEvent`
// and variadic of `SendOption`s.
func (d *Driver) SendInteractive(
	events []*channel.SendInteractiveEvent,
	o ...base.SendOption,
) (*base.Response, error) {
	finalOpts := d.ParseSendOptions(o)
	joinedEventInputs := base.JoinEventInputs(events)

	if d.CurrentPriv != d.DefaultDesiredPriv {
		err := d.AcquirePriv(d.DefaultDesiredPriv)
		if err != nil {
			return nil, err
		}
	}

	return d.Driver.FullSendInteractive(
		events,
		finalOpts.FailedWhenContains,
		finalOpts.TimeoutOps,
		joinedEventInputs,
	)
}
