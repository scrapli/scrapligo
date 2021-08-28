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

	if finalOpts.DesiredPrivilegeLevel == "" {
		finalOpts.DesiredPrivilegeLevel = d.DefaultDesiredPriv
	}

	if d.CurrentPriv != finalOpts.DesiredPrivilegeLevel {
		err := d.AcquirePriv(finalOpts.DesiredPrivilegeLevel)
		if err != nil {
			return nil, err
		}
	}

	return d.Driver.FullSendInteractive(
		events,
		finalOpts.InteractionCompletePatterns,
		finalOpts.FailedWhenContains,
		finalOpts.TimeoutOps,
		joinedEventInputs,
	)
}
