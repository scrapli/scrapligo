package base

import (
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/scrapli/scrapligo/transport"
)

// Option function to set driver options.
type Option func(*Driver) error

// WithPort modify the default (22) port value of a driver.
func WithPort(port int) Option {
	return func(d *Driver) error {
		d.Port = port
		return nil
	}
}

// Auth related options

// WithAuthUsername provide a string username to use for driver authentication.
func WithAuthUsername(username string) Option {
	return func(d *Driver) error {
		d.AuthUsername = username
		return nil
	}
}

// WithAuthPassword provide a string password to use for driver authentication.
func WithAuthPassword(password string) Option {
	return func(d *Driver) error {
		d.AuthPassword = password
		return nil
	}
}

// WithAuthSecondary provide a string "secondary" (or "enable") password to use for driver
// authentication. Only applicable for "network" level drivers.
func WithAuthSecondary(secondary string) Option {
	return func(d *Driver) error {
		d.AuthSecondary = secondary
		return nil
	}
}

// WithAuthPrivateKey provide a string path to a private key to use for driver authentication,
// optionally provide a string to use for passphrase for given private key.
func WithAuthPrivateKey(privateKey string, privateKeyPassphrase ...string) Option {
	pkPassphrase := []string{""}
	if len(privateKeyPassphrase) > 0 {
		pkPassphrase = privateKeyPassphrase
	}

	return func(d *Driver) error {
		d.AuthPrivateKey = privateKey
		d.AuthPrivateKeyPassphrase = strings.Join(pkPassphrase, "")

		return nil
	}
}

// WithAuthBypass provide bool indicating if auth should be "bypassed" -- only applicable for system
// transport.
func WithAuthBypass(bypass bool) Option {
	return func(d *Driver) error {
		d.AuthBypass = bypass
		return nil
	}
}

// WithAuthStrictKey provide bool indicating if strict key checking should be enforced.
func WithAuthStrictKey(stricktKey bool) Option {
	return func(d *Driver) error {
		d.AuthStrictKey = stricktKey
		return nil
	}
}

// SSH file related options

// WithSSHConfigFile provide string path to ssh config file.
func WithSSHConfigFile(sshConfigFile string) Option {
	return func(d *Driver) error {
		d.SSHConfigFile = sshConfigFile
		return nil
	}
}

// WithSSHKnownHostsFile provide string path to ssh known hosts file.
func WithSSHKnownHostsFile(sshKnownHostsFile string) Option {
	return func(d *Driver) error {
		d.SSHKnownHostsFile = sshKnownHostsFile
		return nil
	}
}

// Timeout related options

// WithTimeoutSocket provide duration to use for socket timeout.
func WithTimeoutSocket(timeout time.Duration) Option {
	return func(d *Driver) error {
		d.TimeoutSocket = timeout
		return nil
	}
}

// WithTimeoutTransport provide duration to use for transport timeout.
func WithTimeoutTransport(timeout time.Duration) Option {
	return func(d *Driver) error {
		d.TimeoutTransport = timeout
		return nil
	}
}

// WithTimeoutOps provide duration to use for "operation" timeout.
func WithTimeoutOps(timeout time.Duration) Option {
	return func(d *Driver) error {
		d.TimeoutOps = timeout
		return nil
	}
}

// Comms related options

// WithCommsPromptPattern provide string regex pattern to use for prompt pattern, typically not
// necessary if using a network level driver.
func WithCommsPromptPattern(pattern string) Option {
	return func(d *Driver) error {
		d.CommsPromptPattern = regexp.MustCompile(pattern)
		return nil
	}
}

// WithCommsReturnChar provide string to use as the return character, typically can be left default.
func WithCommsReturnChar(char string) Option {
	return func(d *Driver) error {
		d.CommsReturnChar = char
		return nil
	}
}

// ChannelLog option

// WithChannelLog provide an io.Writer object to write channel log data to.
func WithChannelLog(log io.Writer) Option {
	return func(d *Driver) error {
		d.channelLog = log
		return nil
	}
}

// Transport options

// WithTransportType provide string name of type of transport to use.
func WithTransportType(transportType string) Option {
	var finalTransport string

	switch transportType {
	case transport.SystemTransportName:
		finalTransport = transport.SystemTransportName
	case transport.StandardTransportName:
		finalTransport = transport.StandardTransportName
	default:
		return func(d *Driver) error {
			return transport.ErrUnknownTransport
		}
	}

	return func(d *Driver) error {
		d.TransportType = finalTransport
		return nil
	}
}

// WithTransportPtySize provide pty width/height to use.
func WithTransportPtySize(w, h int) Option {
	return func(d *Driver) error {
		d.transportPtyWidth = w
		d.transportPtyHeight = h

		return nil
	}
}

// Network driver options

// WithFailedWhenContains provide a custom slice of strings to use to check if an output is failed
// -- only applicable to network drivers.
func WithFailedWhenContains(failedWhenContains []string) Option {
	return func(d *Driver) error {
		d.FailedWhenContains = failedWhenContains
		return nil
	}
}

// WithPrivilegeLevels provide custom privilege levels to use -- only applicable to network drivers.
func WithPrivilegeLevels(privilegeLevels map[string]*PrivilegeLevel) Option {
	return func(d *Driver) error {
		d.PrivilegeLevels = privilegeLevels
		return nil
	}
}

// WithDefaultDesiredPriv provide custom default preferred privilege level to use -- only applicable
// to network drivers.
func WithDefaultDesiredPriv(defaultDesiredPriv string) Option {
	return func(d *Driver) error {
		d.DefaultDesiredPriv = defaultDesiredPriv
		return nil
	}
}

// Netconf driver options

// WithNetconfServerEcho provide custom default preferred privilege level to use -- only applicable
// for netconf.
func WithNetconfServerEcho(echo bool) Option {
	return func(d *Driver) error {
		d.NetconfEcho = &echo
		return nil
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
