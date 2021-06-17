package core

import (
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

// NewIOSXRDriver return a driver setup for operation with IOSXR devices.
func NewIOSXRDriver(
	host string,
	options ...base.Option,
) (*network.Driver, error) {
	defaultPrivilegeLevels := map[string]*base.PrivilegeLevel{
		"privilege_exec": {
			Pattern:        `(?im)^[\w.\-@/:]{1,63}#\s?$`,
			Name:           privExecPrivLevel,
			PreviousPriv:   "",
			Deescalate:     "",
			Escalate:       "",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
		"configuration": {
			Pattern:        `(?im)^[\w.\-@/:]{1,63}\(config[\w.\-@/:]{0,32}\)#\s?$`,
			Name:           configPrivLevel,
			PreviousPriv:   privExecPrivLevel,
			Deescalate:     "end",
			Escalate:       "configure terminal",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
		"configuration_exclusive": {
			Pattern:        `(?im)^[\w.\-@/:]{1,63}\(config[\w.\-@/:]{0,32}\)#\s?$`,
			Name:           "configuration_exclusive",
			PreviousPriv:   privExecPrivLevel,
			Deescalate:     "end",
			Escalate:       "configure exclusive",
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
		IOSXROnOpen,
		IOSXROnClose,
		options...)

	if err != nil {
		return nil, err
	}

	d.Augments["abortConfig"] = IOSXRAbortConfig

	return d, nil
}

// IOSXROnOpen default on open callable for IOSXR.
func IOSXROnOpen(d *network.Driver) error {
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

// IOSXROnClose default on close callable for IOSXR.
func IOSXROnClose(d *network.Driver) error {
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

// IOSXRAbortConfig abort IOSXR configuration session.
func IOSXRAbortConfig(d *network.Driver) (*base.Response, error) {
	_, err := d.Channel.SendInput("abort", false, false, -1)

	d.CurrentPriv = privExecPrivLevel

	return nil, err
}
