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
	PatternRe      *regexp.Regexp
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

	NetconfEcho bool
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
