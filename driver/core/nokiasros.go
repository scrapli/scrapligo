package core

import (
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

// NewSROSDriver returns a driver setup for operation with Nokia SR OS devices.
func NewSROSDriver(
	host string,
	options ...base.Option,
) (*network.Driver, error) {
	defaultPrivilegeLevels := map[string]*base.PrivilegeLevel{
		"exec": {
			Pattern:        `(?im)^\[.*\]\n[abcd]:\S+@\S+#\s?$`,
			Name:           execPrivLevel,
			PreviousPriv:   "",
			Deescalate:     "",
			Escalate:       "",
			EscalateAuth:   false,
			EscalatePrompt: ``,
		},
		// configuration privilege level maps to the exclusive config mode on SR OS
		"configuration": {
			Pattern:        `(?im)^\*?\(ex\)\[/?\]\n[abcd]:\S+@\S+#\s?$`,
			Name:           configPrivLevel,
			PreviousPriv:   execPrivLevel,
			Deescalate:     "quit-config",
			Escalate:       "edit-config exclusive",
			EscalateAuth:   false,
			EscalatePrompt: ``,
		},
		"configuration-with-path": {
			Pattern:        `(?im)^\*?\(ex\)\[\S{2,}.+\]\n[abcd]:\S+@\S+#\s?$`,
			Name:           "configuration-with-path",
			PreviousPriv:   configPrivLevel,
			Deescalate:     "exit all",
			Escalate:       "",
			EscalateAuth:   false,
			EscalatePrompt: ``,
		},
	}

	defaultFailedWhenContains := []string{
		"CRITICAL:",
		"MAJOR:",
		"MINOR:",
	}

	const defaultDefaultDesiredPriv = execPrivLevel

	d, err := network.NewNetworkDriver(
		host,
		defaultPrivilegeLevels,
		defaultDefaultDesiredPriv,
		defaultFailedWhenContains,
		SROSOnOpen,
		SROSOnClose,
		options...)

	if err != nil {
		return nil, err
	}

	d.Augments["abortConfig"] = SROSAbortConfig

	return d, nil
}

// SROSOnOpen is a default on open callable for SR OS.
func SROSOnOpen(d *network.Driver) error {
	err := d.AcquirePriv(d.DefaultDesiredPriv)
	if err != nil {
		return err
	}

	if _, err = d.SendCommand("environment console width 512", nil); err != nil {
		return err
	}

	if _, err = d.SendCommand("environment more false", nil); err != nil {
		return err
	}

	_, err = d.SendCommand("environment command-completion space false", nil)

	return err
}

// SROSOnClose is a default on close callable for SR OS.
func SROSOnClose(d *network.Driver) error {
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

// SROSAbortConfig aborts SR OS configuration session.
func SROSAbortConfig(d *network.Driver) (*base.Response, error) {
	if _, err := d.Channel.SendInput("discard /", false, false, -1); err != nil {
		return nil, err
	}

	if _, err := d.Channel.SendInput("exit", false, false, -1); err != nil {
		return nil, err
	}

	_, err := d.Channel.SendInput("quit-config", false, false, -1)

	d.CurrentPriv = privExecPrivLevel

	return nil, err
}
