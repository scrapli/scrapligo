package core

import (
	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/driver/network"
)

// NewSonicDriver return a driver setup for operation with Palo Alto Sonic devices.
func NewSonicDriver(
	host string,
	options ...base.Option,
) (*network.Driver, error) {
	defaultPrivilegeLevels := map[string]*base.PrivilegeLevel{
		"exec": {
			Pattern:        `(?im)^[\w\._-]+@[\w\.\(\)_-]+>\s?`,
			Name:           execPrivLevel,
			PreviousPriv:   "",
			Deescalate:         "disable",
			Escalate:           "enable",
			EscalateAuth:       true,
			EscalatePrompt:     `(?im)^[pP]assword:\s?$`,
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
		SonicOnOpen,
		SonicOnClose,
		options...)

	if err != nil {
		return nil, err
	}

	return d, nil
}

// SonicOnOpen default on open callable for Sonic.
func SonicOnOpen(d *network.Driver) error {
	err := d.AcquirePriv(d.DefaultDesiredPriv)
	if err != nil {
		return err
	}

	return nil
}

// SonicOnClose default on close callable for Sonic.
func SonicOnClose(d *network.Driver) error {
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
