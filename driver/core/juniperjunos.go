package core

import (
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

// NewJUNOSDriver return a driver setup for operation with Junos devices.
func NewJUNOSDriver(
	host string,
	options ...base.Option,
) (*network.Driver, error) {
	defaultPrivilegeLevels := map[string]*base.PrivilegeLevel{
		"exec": {
			Pattern:        `(?im)^({\w+:\d}\n){0,1}[\w\-@()/:]{1,63}>\s?$`,
			Name:           execPrivLevel,
			PreviousPriv:   "",
			Deescalate:     "",
			Escalate:       "",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
		"configuration": {
			Pattern:        `(?im)^({\w+:\d}\[edit\]\n){0,1}[\w\-@()/:]{1,63}#\s?$`,
			Name:           configPrivLevel,
			PreviousPriv:   execPrivLevel,
			Deescalate:     "exit configuration-mode",
			Escalate:       "configure",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
		"configuration_exclusive": {
			Pattern:        `(?im)^({\w+:\d}\[edit\]\n){0,1}[\w\-@()/:]{1,63}#\s?$`,
			Name:           "configuration_exclusive",
			PreviousPriv:   execPrivLevel,
			Deescalate:     "exit configuration-mode",
			Escalate:       "configure exclusive",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
		"configuration_private": {
			Pattern:        `(?im)^({\w+:\d}\[edit\]\n){0,1}[\w\-@()/:]{1,63}#\s?$`,
			Name:           "configuration_private",
			PreviousPriv:   execPrivLevel,
			Deescalate:     "exit configuration-mode",
			Escalate:       "configure exclusive",
			EscalateAuth:   false,
			EscalatePrompt: "",
		},
		"shell": {
			Pattern:            `(?im)^.*[%$]\s?$`,
			PatternNotContains: []string{"root"},
			Name:               "shell",
			PreviousPriv:       execPrivLevel,
			Deescalate:         "exit",
			Escalate:           "start shell",
			EscalateAuth:       false,
			EscalatePrompt:     "",
		},
		"root_shell": {
			Pattern:        `(?im)^.*root@[[:ascii:]]*?:?[[:ascii:]]*?[%#]\s?$`,
			Name:           "root_shell",
			PreviousPriv:   execPrivLevel,
			Deescalate:     "exit",
			Escalate:       "start shell user root",
			EscalateAuth:   true,
			EscalatePrompt: `(?im)^[pP]assword:\s?$`,
		},
	}

	defaultFailedWhenContains := []string{
		"is ambiguous",
		"No valid completions",
		"unknown command",
		"syntax error",
	}

	const defaultDefaultDesiredPriv = execPrivLevel

	d, err := network.NewNetworkDriver(
		host,
		defaultPrivilegeLevels,
		defaultDefaultDesiredPriv,
		defaultFailedWhenContains,
		JUNOSOnOpen,
		JUNOSOnClose,
		options...)

	if err != nil {
		return nil, err
	}

	d.Augments["abortConfig"] = JUNOSAbortConfig

	return d, nil
}

// JUNOSOnOpen default on open callable for Junos.
func JUNOSOnOpen(d *network.Driver) error {
	err := d.AcquirePriv(d.DefaultDesiredPriv)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("set cli screen-length 0", nil)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("set cli screen-width 511", nil)
	if err != nil {
		return err
	}

	_, err = d.SendCommand("set cli complete-on-space off", nil)
	if err != nil {
		return err
	}

	return nil
}

// JUNOSOnClose default on close callable for Junos.
func JUNOSOnClose(d *network.Driver) error {
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

// JUNOSAbortConfig abort Junos configuration session.
func JUNOSAbortConfig(d *network.Driver) (*base.Response, error) {
	_, _ = d.Channel.SendInput("rollback 0", false, false, -1)
	_, err := d.Channel.SendInput("exit", false, false, -1)
	d.CurrentPriv = privExecPrivLevel

	return nil, err
}
