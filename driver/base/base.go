package base

import (
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/transport"
)

// PrivilegeLevel struct defining a single privilege level -- used only for "network" level drivers.
type PrivilegeLevel struct {
	Pattern        string
	Name           string
	PreviousPriv   string
	Deescalate     string
	Escalate       string
	EscalateAuth   bool
	EscalatePrompt string
}

// Driver primary/base driver struct.
type Driver struct {
	Host string
	Port int

	AuthUsername             string
	AuthPassword             string
	AuthSecondary            string
	AuthPrivateKey           string
	AuthPrivateKeyPassphrase string
	AuthStrictKey            bool
	AuthBypass               bool

	SSHConfigFile     string
	SSHKnownHostsFile string

	TimeoutSocket    time.Duration
	TimeoutTransport time.Duration
	TimeoutOps       time.Duration

	CommsPromptPattern *regexp.Regexp
	CommsReturnChar    string

	TransportType      string
	Transport          transport.BaseTransport
	transportPtyWidth  int
	transportPtyHeight int

	Channel    *channel.Channel
	channelLog io.Writer

	FailedWhenContains []string

	PrivilegeLevels    map[string]*PrivilegeLevel
	DefaultDesiredPriv string
}

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
		if d.TransportType == transport.SystemTransportName {
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
		} else if d.TransportType == transport.StandardTransportName {
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

// Open open the connection.
func (d *Driver) Open() error {
	logging.LogDebug(
		d.FormatLogMessage(
			"info",
			fmt.Sprintf("opening connection to '%s' on port '%d'", d.Host, d.Port),
		),
	)

	err := d.Transport.Open()
	if err != nil {
		return err
	}

	if d.TransportType == transport.SystemTransportName && !d.AuthBypass {
		_, err = d.Channel.AuthenticateSSH(d.AuthPassword, d.AuthPrivateKeyPassphrase)
		if err != nil {
			logging.LogError(
				d.FormatLogMessage("error", "authentication failed, connection not opened"),
			)

			return err
		}
	}

	logging.LogDebug(d.FormatLogMessage("info", "connection to device opened successfully"))

	return nil
}

// Close close the connection.
func (d *Driver) Close() error {
	logging.LogDebug(
		d.FormatLogMessage(
			"info",
			fmt.Sprintf("closing connection to '%s' on port '%d'", d.Host, d.Port),
		),
	)

	err := d.Transport.Close()
	if err != nil {
		logging.LogError("failed closing transport")

		return err
	}

	logging.LogDebug(d.FormatLogMessage("info", "connection to device closed successfully"))

	return nil
}

// FormatLogMessage formats log message payload, adding contextual info about the host.
func (d *Driver) FormatLogMessage(level, msg string) string {
	return logging.FormatLogMessage(level, d.Host, d.Port, msg)
}
