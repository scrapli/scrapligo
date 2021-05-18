package core

import (
	"strings"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/driver/network"
)

// NewNXOSDriver return a driver setup for operation with NXOS devices.
func NewNXOSDriver(
	host string,
	options ...base.Option,
) (*network.Driver, error) {
	defaultPrivilegeLevels := map[string]*base.PrivilegeLevel{
		"exec": {
			Pattern:        `(?im)^[a-z0-9.\\-_@()/:]{1,63}>\s?$`,
			Name:           execPrivLevel,
			PreviousPriv:   "",
			Deescalate:     "",
			Escalate:       "",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
		"privilege_exec": {
			Pattern:        `(?im)^[a-z0-9.\-_@/:]{1,63}#\s?$`,
			Name:           privExecPrivLevel,
			PreviousPriv:   execPrivLevel,
			Deescalate:     "disable",
			Escalate:       "enable",
			EscalateAuth:   true,
			EscalatePrompt: `(?im)^(?:enable\s){0,1}password:\s?$`,
		},
		"configuration": {
			Pattern:        `(?im)^[a-z0-9.\-_@/:]{1,63}\(config[a-z0-9.\-@/:\+]{0,32}\)#\s?$`,
			Name:           configPrivLevel,
			PreviousPriv:   privExecPrivLevel,
			Deescalate:     "end",
			Escalate:       "configure terminal",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
	}

	defaultFailedWhenContains := []string{
		"% Ambiguous command",
		"% Incomplete command",
		"% Invalid input detected",
		"% Unknown command",
	}

	const defaultDefaultDesiredPriv = privExecPrivLevel

	d, err := network.NewNetworkDriver(
		host,
		defaultPrivilegeLevels,
		defaultDefaultDesiredPriv,
		defaultFailedWhenContains,
		NXOSOnOpen,
		NXOSOnClose,
		options...)

	if err != nil {
		return nil, err
	}

	d.Augments["abortConfig"] = NXOSAbortConfig

	return d, nil
}

// NXOSOnOpen default on open callable for NXOS.
func NXOSOnOpen(d *network.Driver) error {
	err := d.AcquirePriv(d.DefaultDesiredPriv)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("terminal length 0", nil)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("terminal width 511", nil)
	if err != nil {
		return err
	}

	return nil
}

// NXOSOnClose default on close callable for NXOS.
func NXOSOnClose(d *network.Driver) error {
	err := d.AcquirePriv(d.DefaultDesiredPriv)
	if err != nil {
		return err
	}

	err = d.Channel.Write([]byte("exit"), false)
	if err != nil {
		return err
	}

	err = d.Channel.SendReturn()
	if err != nil {
		return err
	}

	return nil
}

// NXOSAbortConfig abort NXOS configuration session.
func NXOSAbortConfig(d *network.Driver) (*base.Response, error) {
	if strings.Contains(d.CurrentPriv, "config\\-s") {
		_, err := d.Channel.SendInput("abort", false, false, -1)
		d.CurrentPriv = privExecPrivLevel

		return nil, err
	}

	return nil, nil
}
