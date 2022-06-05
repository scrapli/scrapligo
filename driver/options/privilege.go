package options

import (
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/util"
)

// WithPrivilegeLevels allows for setting a custom map of privilege levels for a network driver.
func WithPrivilegeLevels(privilegeLevels map[string]*network.PrivilegeLevel) util.Option {
	return func(o interface{}) error {
		d, ok := o.(*network.Driver)

		if ok {
			d.PrivilegeLevels = privilegeLevels

			return nil
		}

		return util.ErrIgnoredOption
	}
}

// WithDefaultDesiredPriv modifies a network.Driver's default desired privilege level -- that is
// the privilege level at which "show" commands are sent.
func WithDefaultDesiredPriv(s string) util.Option {
	return func(o interface{}) error {
		d, ok := o.(*network.Driver)

		if ok {
			d.DefaultDesiredPriv = s

			return nil
		}

		return util.ErrIgnoredOption
	}
}
