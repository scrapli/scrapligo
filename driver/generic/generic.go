package generic

import (
	"github.com/scrapli/scrapligo/driver/base"
)

// Driver generic driver extending the base driver.
type Driver struct {
	base.Driver
}

// NewGenericDriver return a new Generic Driver.
func NewGenericDriver(
	host string,
	options ...base.Option,
) (*Driver, error) {
	newDriver, err := base.NewDriver(host, options...)

	if err != nil {
		return nil, err
	}

	d := &Driver{
		Driver: *newDriver,
	}

	return d, nil
}

// ParseSendOptions convenience function to parse and set defaults for `SendOption`s.
func (d *Driver) ParseSendOptions(
	o []base.SendOption,
) *base.SendOptions {
	finalOpts := &base.SendOptions{
		StripPrompt:        base.DefaultSendOptionsStripPrompt,
		FailedWhenContains: d.FailedWhenContains,
		StopOnFailed:       base.DefaultSendOptionsStopOnFailed,
		TimeoutOps:         base.DefaultSendOptionsTimeoutOps,
		Eager:              base.DefaultSendOptionsEager,
		// only used with SendConfig(s), thus this should default to "configuration"
		DesiredPrivilegeLevel: "configuration",
	}

	if len(o) > 0 && o[0] != nil {
		for _, option := range o {
			option(finalOpts)
		}
	}

	return finalOpts
}
