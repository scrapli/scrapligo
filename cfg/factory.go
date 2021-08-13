package cfg

import (
	"errors"

	"github.com/scrapli/scrapligo/driver/network"
)

// ErrUnknownCfgPlatform raised when user provides an unknown cfg platform... duh.
var ErrUnknownCfgPlatform = errors.New("unknown cfg platform provided")

// SupportedPlatforms pseudo constant providing slice of all core cfg platform types.
func SupportedPlatforms() []string {
	return []string{"cisco_iosxe",
		"cisco_iosxr",
		"cisco_nxos",
		"arista_eos",
		"juniper_junos",
	}
}

// NewCfgDriver return a new cfg driver for a given platform.
func NewCfgDriver(
	conn *network.Driver,
	platform string,
	options ...Option,
) (*Cfg, error) {
	switch platform {
	case "cisco_iosxe":
		return NewIOSXECfg(conn, options...)
	case "cisco_iosxr":
		return NewIOSXRCfg(conn, options...)
	case "cisco_nxos":
		return NewNXOSCfg(conn, options...)
	case "arista_eos":
		return NewEOSCfg(conn, options...)
	case "juniper_junos":
		return NewJUNOSCfg(conn, options...)
	}

	return nil, ErrUnknownCfgPlatform
}
