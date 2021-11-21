package core

import (
	"errors"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/driver/network"
)

// ErrUnknownPlatform raised when user provides an unknown platform... duh.
var ErrUnknownPlatform = errors.New("unknown platform provided")

// SupportedPlatforms pseudo constant providing slice of all core platform types.
func SupportedPlatforms() []string {
	return []string{"cisco_iosxe",
		"cisco_iosxr",
		"cisco_nxos",
		"arista_eos",
		"juniper_junos",
		"nokia_sros",
		"nokia_sros_classic",
		"paloalto_panos",
	}
}

// NewCoreDriver return a new core driver for a given platform.
var NewCoreDriver = newCoreDriver //nolint:gochecknoglobals

func newCoreDriver(
	host,
	platform string,
	options ...base.Option,
) (*network.Driver, error) {
	switch platform {
	case "cisco_iosxe":
		return NewIOSXEDriver(host, options...)
	case "cisco_iosxr":
		return NewIOSXRDriver(host, options...)
	case "cisco_nxos":
		return NewNXOSDriver(host, options...)
	case "arista_eos":
		return NewEOSDriver(host, options...)
	case "juniper_junos":
		return NewJUNOSDriver(host, options...)
	case "nokia_sros":
		return NewSROSDriver(host, options...)
	case "nokia_sros_classic":
		return NewSROSClassicDriver(host, options...)
	case "paloalto_panos":
		return NewPanOSDriver(host, options...)
	}

	return nil, ErrUnknownPlatform
}
