package core

import (
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

// NewIOSXEDriver return a driver setup for operation with IOSXE devices.
func NewIOSXEDriver(
	host string,
	options ...base.Option,
) (*network.Driver, error) {
	defaultPrivilegeLevels := map[string]*base.PrivilegeLevel{
		"exec": {
			Pattern:        `(?im)^[\w.\-@/:]{1,63}>$`,
			Name:           execPrivLevel,
			PreviousPriv:   "",
			Deescalate:     "",
			Escalate:       "",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
		"privilege_exec": {
			Pattern:        `(?im)^[\w.\-@/:]{1,63}#$`,
			Name:           privExecPrivLevel,
			PreviousPriv:   execPrivLevel,
			Deescalate:     "disable",
			Escalate:       "enable",
			EscalateAuth:   true,
			EscalatePrompt: `(?im)^(?:enable\s){0,1}password:\s?$`,
		},
		"configuration": {
			Pattern:            `(?im)^[\w.\-@/:]{1,63}\([\w.\-@/:+]{0,32}\)#$`,
			PatternNotContains: []string{"tcl)"},
			Name:               configPrivLevel,
			PreviousPriv:       privExecPrivLevel,
			Deescalate:         "end",
			Escalate:           "configure terminal",
			EscalateAuth:       false,
			EscalatePrompt:     "",
		},
		"tclsh": {
			Pattern:        `(?im)^([\w.\-@/+>:]+\(tcl\)[>#]|\+>)$`,
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
		IOSXEOnOpen,
		IOSXEOnClose,
		options...)

	if err != nil {
		return nil, err
	}

	return d, nil
}

// IOSXEOnOpen default on open callable for IOSXE.
func IOSXEOnOpen(d *network.Driver) error {
	err := d.AcquirePriv(d.DefaultDesiredPriv)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("terminal length 0", nil)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("terminal width 512", nil)
	if err != nil {
		return err
	}

	return nil
}

// IOSXEOnClose default on close callable for IOSXE.
func IOSXEOnClose(d *network.Driver) error {
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
