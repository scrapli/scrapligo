package core

import (
	"strings"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/driver/network"
)

// NewEOSDriver return a driver setup for operation with EOS devices.
func NewEOSDriver(
	host string,
	options ...base.Option,
) (*network.Driver, error) {
	defaultPrivilegeLevels := map[string]*base.PrivilegeLevel{
		"exec": {
			Pattern:        `(?im)^[\w.\-@()/: ]{1,63}>\s?$`,
			Name:           execPrivLevel,
			PreviousPriv:   "",
			Deescalate:     "",
			Escalate:       "",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
		"privilege_exec": {
			Pattern:            `(?im)^[\w.\-@()/: ]{1,63}#\s?$`,
			PatternNotContains: []string{"(config"},
			Name:               privExecPrivLevel,
			PreviousPriv:       execPrivLevel,
			Deescalate:         "disable",
			Escalate:           "enable",
			EscalateAuth:       true,
			EscalatePrompt:     `(?im)^[pP]assword:\s?$`,
		},
		"configuration": {
			Pattern:            `(?im)^[\w.\-@()/: ]{1,63}\(config[\w.\-@/:]{0,32}\)#\s?$`,
			PatternNotContains: []string{"(config-s-"},
			Name:               configPrivLevel,
			PreviousPriv:       privExecPrivLevel,
			Deescalate:         "end",
			Escalate:           "configure terminal",
			EscalateAuth:       false,
			EscalatePrompt:     "",
		},
	}
	defaultFailedWhenContains := []string{
		"% Ambiguous command",
		"% Error",
		"% Incomplete command",
		"% Invalid input",
		"% Cannot commit",
		"% Unavailable command",
	}

	const defaultDefaultDesiredPriv = privExecPrivLevel

	d, err := network.NewNetworkDriver(
		host,
		defaultPrivilegeLevels,
		defaultDefaultDesiredPriv,
		defaultFailedWhenContains,
		EOSOnOpen,
		EOSOnClose,
		options...)

	if err != nil {
		return nil, err
	}

	d.Augments["abortConfig"] = EOSAbortConfig

	return d, nil
}

// EOSOnOpen default on open callable for EOS.
func EOSOnOpen(d *network.Driver) error {
	err := d.AcquirePriv(d.DefaultDesiredPriv)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("terminal length 0", nil)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("terminal width 32767", nil)
	if err != nil {
		return err
	}

	return nil
}

// EOSOnClose default on close callable for EOS.
func EOSOnClose(d *network.Driver) error {
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

// EOSAbortConfig abort EOS configuration session.
func EOSAbortConfig(d *network.Driver) (*base.Response, error) {
	if strings.Contains(d.CurrentPriv, "config\\-s") {
		_, err := d.Channel.SendInput("abort", false, false, -1)
		d.CurrentPriv = privExecPrivLevel

		return nil, err
	}

	return nil, nil
}
