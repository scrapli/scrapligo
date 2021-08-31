package base

import (
	"fmt"
	"regexp"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/transport"
)

// PrivilegeLevel struct defining a single privilege level -- used only for "network" level drivers.
type PrivilegeLevel struct {
	Pattern            string
	PatternRe          *regexp.Regexp
	PatternNotContains []string
	Name               string
	PreviousPriv       string
	Deescalate         string
	Escalate           string
	EscalateAuth       bool
	EscalatePrompt     string
}

// Driver primary/base driver struct.
type Driver struct {
	Host string

	AuthUsername             string
	AuthPassword             string
	AuthSecondary            string
	AuthPrivateKey           string
	AuthPrivateKeyPassphrase string
	AuthStrictKey            bool
	AuthBypass               bool

	SSHConfigFile     string
	SSHKnownHostsFile string

	TransportType string
	Transport     *transport.Transport

	Channel *channel.Channel

	FailedWhenContains []string

	PrivilegeLevels    map[string]*PrivilegeLevel
	DefaultDesiredPriv string

	NetconfEcho *bool
}

// Open opens the connection.
func (d *Driver) Open() error {
	logging.LogDebug(
		d.FormatLogMessage(
			"info",
			fmt.Sprintf(
				"opening connection to '%s' on port '%d'",
				d.Host,
				d.Transport.BaseTransportArgs.Port,
			),
		),
	)

	err := d.Transport.Open()
	if err != nil {
		return err
	}

	var authErr error
	if d.TransportType == transport.SystemTransportName && !d.AuthBypass {
		_, authErr = d.Channel.AuthenticateSSH(d.AuthPassword, d.AuthPrivateKeyPassphrase)
	} else if d.TransportType == transport.TelnetTransportName && !d.AuthBypass {
		_, authErr = d.Channel.AuthenticateTelnet(d.AuthUsername, d.AuthPassword)
	}

	if authErr != nil {
		logging.LogError(
			d.FormatLogMessage("error", "authentication failed, connection not opened"),
		)

		return err
	}

	logging.LogDebug(d.FormatLogMessage("info", "connection to device opened successfully"))

	return nil
}

// Close closes the connection.
func (d *Driver) Close() error {
	logging.LogDebug(
		d.FormatLogMessage(
			"info",
			fmt.Sprintf(
				"closing connection to '%s' on port '%d'",
				d.Host,
				d.Transport.BaseTransportArgs.Port,
			),
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
	return logging.FormatLogMessage(level, d.Host, d.Transport.BaseTransportArgs.Port, msg)
}
