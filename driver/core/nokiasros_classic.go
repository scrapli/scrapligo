package core

import (
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

// NewSROSClassicDriver returns a driver setup for operation
// with Nokia SR OS devices running in classic configuration mode.
func NewSROSClassicDriver(
	host string,
	options ...base.Option,
) (*network.Driver, error) {
	defaultPrivilegeLevels := map[string]*base.PrivilegeLevel{
		"configuration": {
			Pattern:        `(?im)^\*?[abcd]:\S+#\s*$`,
			Name:           configPrivLevel,
			PreviousPriv:   "",
			Deescalate:     "",
			Escalate:       "",
			EscalateAuth:   false,
			EscalatePrompt: ``,
		},
	}

	defaultFailedWhenContains := []string{
		"CRITICAL:",
		"MAJOR:",
		"MINOR:",
		"Error:",
	}

	const defaultDefaultDesiredPriv = configPrivLevel

	d, err := network.NewNetworkDriver(
		host,
		defaultPrivilegeLevels,
		defaultDefaultDesiredPriv,
		defaultFailedWhenContains,
		SROSClassicOnOpen,
		SROSClassicOnClose,
		options...)

	if err != nil {
		return nil, err
	}

	return d, nil
}

// SROSClassicOnOpen is a default on open callable for SR OS classic.
func SROSClassicOnOpen(d *network.Driver) error {
	err := d.AcquirePriv(d.DefaultDesiredPriv)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("environment no more", nil)

	return err
}

// SROSClassicOnClose is a default on close callable for SR OS classic.
func SROSClassicOnClose(d *network.Driver) error {
	err := d.AcquirePriv(d.DefaultDesiredPriv)
	if err != nil {
		return err
	}

	err = d.Channel.Write([]byte("logout"), false)
	if err != nil {
		return err
	}

	err = d.Channel.SendReturn()
	if err != nil {
		return err
	}

	return nil
}
