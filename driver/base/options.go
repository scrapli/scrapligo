package base

import (
	"errors"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/channel"

	"github.com/scrapli/scrapligo/transport"
)

var ErrIgnoredOption = errors.New("option ignored, for different instance type")

// Option function to set driver options.
type Option func(interface{}) error

// WithAuthUsername provide a string username to use for driver authentication.
func WithAuthUsername(username string) Option {
	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			d.AuthUsername = username
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithAuthPassword provide a string password to use for driver authentication.
func WithAuthPassword(password string) Option {
	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			d.AuthPassword = password
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithAuthSecondary provide a string "secondary" (or "enable") password to use for driver
// authentication. Only applicable for "network" level drivers.
func WithAuthSecondary(secondary string) Option {
	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			d.AuthSecondary = secondary
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithAuthPrivateKey provide a string path to a private key to use for driver authentication,
// optionally provide a string to use for passphrase for given private key.
func WithAuthPrivateKey(privateKey string, privateKeyPassphrase ...string) Option {
	pkPassphrase := []string{""}
	if len(privateKeyPassphrase) > 0 {
		pkPassphrase = privateKeyPassphrase
	}

	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			d.AuthPrivateKey = privateKey
			d.AuthPrivateKeyPassphrase = strings.Join(pkPassphrase, "")

			return nil
		}

		return ErrIgnoredOption
	}
}

// WithAuthBypass provide bool indicating if auth should be "bypassed" -- only applicable for system
// and telnet transports.
func WithAuthBypass(bypass bool) Option {
	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			d.AuthBypass = bypass
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithAuthStrictKey provide bool indicating if strict key checking should be enforced.
func WithAuthStrictKey(stricktKey bool) Option {
	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			d.AuthStrictKey = stricktKey
			return nil
		}

		return ErrIgnoredOption
	}
}

// SSH file related options

// WithSSHConfigFile provide string path to ssh config file.
func WithSSHConfigFile(sshConfigFile string) Option {
	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			resolvedSSHConfigFile, err := util.ResolveFilePath(sshConfigFile)
			if err != nil {
				return err
			}

			d.SSHConfigFile = resolvedSSHConfigFile

			return nil
		}

		return ErrIgnoredOption
	}
}

// WithSSHKnownHostsFile provide string path to ssh known hosts file.
func WithSSHKnownHostsFile(sshKnownHostsFile string) Option {
	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			resolvedSSHKnownHostsFile, err := util.ResolveFilePath(sshKnownHostsFile)
			if err != nil {
				return err
			}

			d.SSHKnownHostsFile = resolvedSSHKnownHostsFile

			return nil
		}

		return ErrIgnoredOption
	}
}

// Network driver options

// WithFailedWhenContains provide a custom slice of strings to use to check if an output is failed
// -- only applicable to network drivers.
func WithFailedWhenContains(failedWhenContains []string) Option {
	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			d.FailedWhenContains = failedWhenContains
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithPrivilegeLevels provide custom privilege levels to use -- only applicable to network drivers.
func WithPrivilegeLevels(privilegeLevels map[string]*PrivilegeLevel) Option {
	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			d.PrivilegeLevels = privilegeLevels
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithDefaultDesiredPriv provide custom default preferred privilege level to use -- only applicable
// to network drivers.
func WithDefaultDesiredPriv(defaultDesiredPriv string) Option {
	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			d.DefaultDesiredPriv = defaultDesiredPriv
			return nil
		}

		return ErrIgnoredOption
	}
}

// Netconf driver options

// WithNetconfServerEcho provide custom default preferred privilege level to use -- only applicable
// for netconf.
func WithNetconfServerEcho(echo bool) Option {
	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			d.NetconfEcho = &echo
			return nil
		}

		return ErrIgnoredOption
	}
}

// Send command/config options

const (
	// DefaultSendOptionsStripPrompt default to stripping prompt.
	DefaultSendOptionsStripPrompt = true
	// DefaultSendOptionsStopOnFailed default to *not* stopping on failures.
	DefaultSendOptionsStopOnFailed = false
	// DefaultSendOptionsTimeoutOps default to relying on the drivers timeout ops attribute.
	DefaultSendOptionsTimeoutOps = -1.0
	// DefaultSendOptionsEager default to *not* eager mode.
	DefaultSendOptionsEager = false
)

// SendOptions struct for send operation options.
type SendOptions struct {
	StripPrompt                 bool
	FailedWhenContains          []string
	StopOnFailed                bool
	TimeoutOps                  time.Duration
	Eager                       bool
	DesiredPrivilegeLevel       string
	InteractionCompletePatterns []string
}

// SendOption func to set send options.
type SendOption func(*SendOptions)

// WithSendStripPrompt bool indicating if you would like the hostname/device prompt stripped out of
// output from a send operation.
func WithSendStripPrompt(stripPrompt bool) SendOption {
	return func(o *SendOptions) {
		o.StripPrompt = stripPrompt
	}
}

// WithSendFailedWhenContains slice of strings that overrides the drivers `FailedWhenContains` list
// for a given send operation.
func WithSendFailedWhenContains(failedWhenContains []string) SendOption {
	return func(o *SendOptions) {
		o.FailedWhenContains = failedWhenContains
	}
}

// WithSendStopOnFailed bool indicating if multi command/config operations should stop at first sign
// of failure (based on FailedWhenContains list).
func WithSendStopOnFailed(stopOnFailed bool) SendOption {
	return func(o *SendOptions) {
		o.StopOnFailed = stopOnFailed
	}
}

// WithSendTimeoutOps duration to use for timeout of a given send operation.
func WithSendTimeoutOps(timeoutOps time.Duration) SendOption {
	return func(o *SendOptions) {
		o.TimeoutOps = timeoutOps
	}
}

// WithSendEager bool indicating if send operation should operate in `eager` mode -- generally only
// used for netconf operations.
func WithSendEager(eager bool) SendOption {
	return func(o *SendOptions) {
		o.Eager = eager
	}
}

// WithDesiredPrivilegeLevel provide a desired privilege level for the send operation to work in.
func WithDesiredPrivilegeLevel(privilegeLevel string) SendOption {
	// ignored for command(s) operations, only applicable for interactive/config operations
	return func(o *SendOptions) {
		o.DesiredPrivilegeLevel = privilegeLevel
	}
}

// WithInteractionCompletePatterns provide a list of patterns which, when seen, indicate a
// `SendInteractive` "session" is complete. Only used for `SendInteractive`, otherwise ignored.
func WithInteractionCompletePatterns(interactionCompletePatterns []string) SendOption {
	// ignored for command(s) operations, only applicable for interactive/config operations
	return func(o *SendOptions) {
		o.InteractionCompletePatterns = interactionCompletePatterns
	}
}

// WithPort modifies the default (22) port value of a driver/transport.
func WithPort(port int) Option {
	return func(o interface{}) error {
		t, ok := o.(*transport.Transport)

		if ok {
			t.BaseTransportArgs.Port = port
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithTimeoutSocket provide duration to use for socket timeout.
func WithTimeoutSocket(timeout time.Duration) Option {
	return func(o interface{}) error {
		t, ok := o.(*transport.Transport)

		if ok {
			t.BaseTransportArgs.TimeoutSocket = timeout
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithTimeoutTransport provide duration to use for transport timeout.
func WithTimeoutTransport(timeout time.Duration) Option {
	return func(o interface{}) error {
		t, ok := o.(*transport.Transport)

		if ok {
			t.BaseTransportArgs.TimeoutTransport = timeout
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithTransportType provide string name of type of transport to use.
func WithTransportType(transportType string) Option {
	var finalTransport string

	switch transportType {
	case transport.SystemTransportName:
		finalTransport = transport.SystemTransportName
	case transport.StandardTransportName:
		finalTransport = transport.StandardTransportName
	case transport.TelnetTransportName:
		finalTransport = transport.TelnetTransportName
	default:
		return func(o interface{}) error {
			return transport.ErrUnknownTransport
		}
	}

	return func(o interface{}) error {
		d, ok := o.(*Driver)

		if ok {
			d.TransportType = finalTransport
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithTransportPtySize provide pty width/height to use.
func WithTransportPtySize(w, h int) Option {
	return func(o interface{}) error {
		t, ok := o.(*transport.Transport)

		if ok {
			t.BaseTransportArgs.PtyWidth = w
			t.BaseTransportArgs.PtyHeight = h

			return nil
		}

		return ErrIgnoredOption
	}
}

// WithTimeoutOps provide duration to use for "operation" timeout.
func WithTimeoutOps(timeout time.Duration) Option {
	return func(o interface{}) error {
		c, ok := o.(*channel.Channel)

		if ok {
			c.TimeoutOps = timeout
			return nil
		}

		return ErrIgnoredOption
	}
}

// Comms related options

// WithCommsPromptPattern provide string regex pattern to use for prompt pattern, typically not
// necessary if using a network level driver.
func WithCommsPromptPattern(pattern string) Option {
	return func(o interface{}) error {
		c, ok := o.(*channel.Channel)

		if ok {
			c.CommsPromptPattern = regexp.MustCompile(pattern)
			return nil
		}

		return ErrIgnoredOption
	}
}

// WithCommsReturnChar provide string to use as the return character, typically can be left default.
func WithCommsReturnChar(char string) Option {
	return func(o interface{}) error {
		c, ok := o.(*channel.Channel)

		if ok {
			c.CommsReturnChar = char
			return nil
		}

		return ErrIgnoredOption
	}
}

// ChannelLog option

// WithChannelLog provide an io.Writer object to write channel log data to.
func WithChannelLog(log io.Writer) Option {
	return func(o interface{}) error {
		c, ok := o.(*channel.Channel)

		if ok {
			c.ChannelLog = log
			return nil
		}

		return ErrIgnoredOption
	}
}
