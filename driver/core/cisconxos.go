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
			Pattern:        `(?im)^[\w.\-]{1,63}>\s?$`,
			Name:           execPrivLevel,
			PreviousPriv:   "",
			Deescalate:     "",
			Escalate:       "",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
		"privilege_exec": {
			Pattern:            `(?im)^[\w.\-]{1,63}#\s?$`,
			PatternNotContains: []string{"-tcl"},
			Name:               privExecPrivLevel,
			PreviousPriv:       execPrivLevel,
			Deescalate:         "disable",
			Escalate:           "enable",
			EscalateAuth:       true,
			EscalatePrompt:     `(?im)^[pP]assword:\s?$`,
		},
		"configuration": {
			Pattern:            `(?im)^[\w.\-]{1,63}\(config[\w.\-@/:]{0,32}\)#\s?$`,
			PatternNotContains: []string{"config-tcl", "config-s"},
			Name:               configPrivLevel,
			PreviousPriv:       privExecPrivLevel,
			Deescalate:         "end",
			Escalate:           "configure terminal",
			EscalateAuth:       false,
			EscalatePrompt:     "",
		},
		"tclsh": {
			Pattern: `(?im)(^[\w.\-@/:]{1,63}\-tcl#\s?$)|` +
				`(^[\w.\-@/:]{1,63}\(config\-tcl\)#\s?$)|(^>\s?$)`,
			Name:           "tclsh",
			PreviousPriv:   privExecPrivLevel,
			Deescalate:     "tclquit",
			Escalate:       "tclsh",
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
