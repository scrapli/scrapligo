package network

import (
	"errors"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/channel"
)

// SendConfigs send configurations to the device.
func (d *Driver) SendConfigs(c []string, o ...base.SendOption) (*base.MultiResponse, error) {
	finalOpts := d.ParseSendOptions(o)

	if finalOpts.DesiredPrivilegeLevel == "" {
		finalOpts.DesiredPrivilegeLevel = "configuration"
	}

	if d.CurrentPriv != finalOpts.DesiredPrivilegeLevel {
		err := d.AcquirePriv(finalOpts.DesiredPrivilegeLevel)
		if err != nil {
			return nil, err
		}
	}

	m, err := d.Driver.FullSendCommands(c,
		finalOpts.FailedWhenContains,
		finalOpts.StripPrompt,
		finalOpts.StopOnFailed,
		finalOpts.Eager,
		finalOpts.TimeoutOps,
	)

	if err != nil && !errors.Is(err, channel.ErrChannelTimeout) {
		// if we encountered an error we *probably* cant abort anyway unless its a timeout error
		// if its a timeout error we can at least try to keep going on, otherwise lets bail here
		return m, err
	}

	if finalOpts.StopOnFailed && m.Failed != nil {
		if f, ok := d.Augments["abortConfig"]; ok {
			_, err = f(d)
		}
	}

	return m, err
}

// SendConfigsFromFile send configurations from a file to the device.
func (d *Driver) SendConfigsFromFile(
	f string,
	o ...base.SendOption,
) (*base.MultiResponse, error) {
	c, err := util.LoadFileLines(f)
	if err != nil {
		return nil, err
	}

	return d.SendConfigs(c, o...)
}
