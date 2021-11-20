package core

import (
	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/driver/network"
)

// NewPanOSDriver return a driver setup for operation with Palo Alto PanOS devices.
func NewPanOSDriver(
	host string,
	options ...base.Option,
) (*network.Driver, error) {
	defaultPrivilegeLevels := map[string]*base.PrivilegeLevel{
		"exec": {
			Pattern:        `(?im)^[\w\._-]+@[\w\.\(\)_-]+>\s?`,
			Name:           execPrivLevel,
			PreviousPriv:   "",
			Deescalate:     "",
			Escalate:       "",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
		"configuration": {
			Pattern:        `(?im)^[\w\._-]+@[\w\.\(\)_-]+#\s?$`,
			Name:           configPrivLevel,
			PreviousPriv:   execPrivLevel,
			Deescalate:     "exit",
			Escalate:       "configure",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
	}
	defaultFailedWhenContains := []string{
		"Unknown command:",
		"Invalid Syntax.",
		"Validation Error:",
	}

	const defaultDefaultDesiredPriv = execPrivLevel

	d, err := network.NewNetworkDriver(
		host,
		defaultPrivilegeLevels,
		defaultDefaultDesiredPriv,
		defaultFailedWhenContains,
		PanOSOnOpen,
		PanOSOnClose,
		options...)

	if err != nil {
		return nil, err
	}

	return d, nil
}

// PanOSOnOpen default on open callable for PanOS.
func PanOSOnOpen(d *network.Driver) error {
	err := d.AcquirePriv(d.DefaultDesiredPriv)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("set cli scripting-mode on", nil)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("set cli pager off", nil)
	if err != nil {
		return err
	}

	return nil
}

// PanOSOnClose default on close callable for PanOS.
func PanOSOnClose(d *network.Driver) error {
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
