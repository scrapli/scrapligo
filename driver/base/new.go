// +build !windows

package base

import (
	"regexp"
	"time"

	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/transport"
)

// NewDriver create a new instance of `Driver`, accepts a host and variadic of options to modify
// the driver behavior.
func NewDriver(
	host string,
	options ...Option,
) (*Driver, error) {
	d := &Driver{
		Host:               host,
		Port:               22,
		AuthStrictKey:      true,
		TimeoutSocket:      30 * time.Second,
		TimeoutTransport:   45 * time.Second,
		TimeoutOps:         60 * time.Second,
		CommsPromptPattern: regexp.MustCompile(`(?im)^[a-z0-9.\-@()/:]{1,48}[#>$]\s*$`),
		CommsReturnChar:    "\n",
		TransportType:      transport.SystemTransportName,
		transportPtyHeight: 80,
		transportPtyWidth:  256,
		FailedWhenContains: []string{},
		PrivilegeLevels:    map[string]*PrivilegeLevel{},
		DefaultDesiredPriv: "",
	}

	for _, option := range options {
		if err := option(d); err != nil {
			return nil, err
		}
	}

	baseTransportArgs := &transport.BaseTransportArgs{
		Host:             d.Host,
		Port:             d.Port,
		AuthUsername:     d.AuthUsername,
		TimeoutSocket:    &d.TimeoutSocket,
		TimeoutTransport: &d.TimeoutTransport,
		PtyHeight:        d.transportPtyHeight,
		PtyWidth:         d.transportPtyWidth,
	}

	if d.Transport == nil {
		switch d.TransportType {
		case transport.SystemTransportName:
			systemTransportArgs := &transport.SystemTransportArgs{
				AuthPrivateKey:    d.AuthPrivateKey,
				AuthStrictKey:     d.AuthStrictKey,
				SSHConfigFile:     d.SSHConfigFile,
				SSHKnownHostsFile: d.SSHKnownHostsFile,
			}
			t := &transport.System{
				BaseTransportArgs:   baseTransportArgs,
				SystemTransportArgs: systemTransportArgs,
			}
			d.Transport = t
		case transport.StandardTransportName:
			standardTransportArgs := &transport.StandardTransportArgs{
				AuthPassword:      d.AuthPassword,
				AuthPrivateKey:    d.AuthPrivateKey,
				AuthStrictKey:     d.AuthStrictKey,
				SSHConfigFile:     d.SSHConfigFile,
				SSHKnownHostsFile: d.SSHKnownHostsFile,
			}
			t := &transport.Standard{
				BaseTransportArgs:     baseTransportArgs,
				StandardTransportArgs: standardTransportArgs,
			}
			d.Transport = t
		default:
			return nil, transport.ErrUnknownTransport
		}
	}

	c := &channel.Channel{
		CommsPromptPattern: d.CommsPromptPattern,
		CommsReturnChar:    &d.CommsReturnChar,
		TimeoutOps:         &d.TimeoutOps,
		Transport:          d.Transport,
		Host:               d.Host,
		Port:               d.Port,
		ChannelLog:         d.channelLog,
	}

	d.Channel = c

	return d, nil
}
