package driver

import (
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
)

// Option is a type used for functional options for the Driver object's options.
type Option func(d *Driver) error

// TransportKind is an enum(ish) representing the kind of transport a Driver should use.
type TransportKind string

const (
	// TransportKindBin represents the "bin" transport -- the default transport that is a wrapper
	// around /bin/ssh.
	TransportKindBin TransportKind = "bin"
	// TransportKindSSH2 represents the "ssh2" transport -- the transport using libssh2.
	TransportKindSSH2 TransportKind = "ssh2"
	// TransportKindTelnet represents the "telnet" transport.
	TransportKindTelnet TransportKind = "telnet"
	// TransportKindTest represents the "Test" transport that is used for integration testing.
	TransportKindTest TransportKind = "test"
)

const (
	// DefaultSSHPort is the default port used for SSH operations.
	DefaultSSHPort uint16 = 22
	// DefaultTelnetPort is the default port used for telnet operations.
	DefaultTelnetPort uint16 = 23
)

func newOptions() options {
	return options{
		definitionVariant: "default",
		loggerCallback:    nil,
		transportKind:     TransportKindBin,
		port:              nil,
		auth: authOptions{
			lookupMap: make(map[string]string),
		},
	}
}

type options struct {
	// the loaded string is set in NewDriver, not via an option as the name/file of the definition
	// is a required argument.
	definitionString  string
	definitionVariant string

	loggerCallback func(level uint8, message *string)

	port *uint16

	transportKind TransportKind

	session   sessionOptions
	auth      authOptions
	transport transportOptions
}

func (o *options) apply(driverPtr uintptr, m *scrapligoffi.Mapping) error {
	err := o.session.apply(driverPtr, m)
	if err != nil {
		return err
	}

	err = o.auth.apply(driverPtr, m)
	if err != nil {
		return err
	}

	switch o.transportKind {
	case TransportKindBin:
		err = o.transport.bin.apply(driverPtr, m)
		if err != nil {
			return err
		}
	case TransportKindSSH2:
		err = o.transport.ssh2.apply(driverPtr, m)
		if err != nil {
			return err
		}
	case TransportKindTelnet:
	case TransportKindTest:
	}

	return nil
}

type sessionOptions struct {
	readSize               *uint64
	readDelayMinNs         *uint64
	readDelayMaxNs         *uint64
	readDelayBackoffFactor *uint8
	returnChar             string

	operationTimeoutNs      *uint64
	operationMaxSearchDepth *uint64
}

func (o *sessionOptions) apply( //nolint: gocyclo
	driverPtr uintptr,
	m *scrapligoffi.Mapping,
) error {
	if o.readSize != nil {
		rc := m.Options.Session.SetReadSize(driverPtr, *o.readSize)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting read size option", nil)
		}
	}

	if o.readDelayMinNs != nil {
		rc := m.Options.Session.SetReadDelayMinNs(driverPtr, *o.readDelayMinNs)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting min read delay option", nil)
		}
	}

	if o.readDelayMaxNs != nil {
		rc := m.Options.Session.SetReadDelayMaxNs(driverPtr, *o.readDelayMaxNs)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting max read delay option", nil)
		}
	}

	if o.readDelayBackoffFactor != nil {
		rc := m.Options.Session.SetReadDelayBackoffFactor(driverPtr, *o.readDelayBackoffFactor)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting read delay backoff factor option",
				nil,
			)
		}
	}

	if o.returnChar != "" {
		rc := m.Options.Session.SetReturnChar(driverPtr, o.returnChar)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting return char option", nil)
		}
	}

	if o.operationTimeoutNs != nil {
		rc := m.Options.Session.SetOperationTimeoutNs(driverPtr, *o.operationTimeoutNs)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting operation timeout option", nil)
		}
	} else {
		// if user does not provide a timeout we assume they want to govern all timeouts via context
		// cancellation
		rc := m.Options.Session.SetOperationTimeoutNs(driverPtr, 0)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting operation timeout option", nil)
		}
	}

	if o.operationMaxSearchDepth != nil {
		rc := m.Options.Session.SetOperationMaxSearchDepth(driverPtr, *o.operationMaxSearchDepth)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting operation search depth option",
				nil,
			)
		}
	}

	return nil
}

type authOptions struct {
	username string
	password string

	privateKeyPath       string
	privateKeyPassphrase string

	lookupMap map[string]string

	inSessionAuthBypass bool

	usernamePattern   string
	passwordPattern   string
	passphrasePattern string
}

func (o *authOptions) apply(driverPtr uintptr, m *scrapligoffi.Mapping) error { //nolint: gocyclo
	if o.username != "" {
		rc := m.Options.Auth.SetUsername(driverPtr, o.username)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting username option", nil)
		}
	}

	if o.password != "" {
		rc := m.Options.Auth.SetPassword(driverPtr, o.password)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting password option", nil)
		}
	}

	if o.privateKeyPath != "" {
		rc := m.Options.Auth.SetPrivateKeyPath(driverPtr, o.privateKeyPath)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting private key path option", nil)
		}
	}

	if o.privateKeyPassphrase != "" {
		rc := m.Options.Auth.SetPrivateKeyPassphrase(driverPtr, o.privateKeyPassphrase)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting private key passphrase option",
				nil,
			)
		}
	}

	for k, v := range o.lookupMap {
		rc := m.Options.Auth.SetDriverOptionAuthLookupKeyValue(driverPtr, k, v)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting lookup map option",
				nil,
			)
		}
	}

	if o.inSessionAuthBypass {
		rc := m.Options.Auth.SetInSessionAuthBypass(driverPtr)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting in session auth bypass option",
				nil,
			)
		}
	}

	if o.usernamePattern != "" {
		rc := m.Options.Auth.SetUsernamePattern(driverPtr, o.usernamePattern)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting username pattern option", nil)
		}
	}

	if o.passwordPattern != "" {
		rc := m.Options.Auth.SetPasswordPattern(driverPtr, o.passwordPattern)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting password pattern option", nil)
		}
	}

	if o.passphrasePattern != "" {
		rc := m.Options.Auth.SetPassphrasePattern(driverPtr, o.passphrasePattern)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting passphrase pattern option", nil)
		}
	}

	return nil
}

type transportOptions struct {
	bin  transportBinOptions
	ssh2 transportSSH2Options
}

type transportBinOptions struct {
	bin              string
	extraOpenArgs    string
	overrideOpenArgs string
	sshConfigPath    string
	knownHostsPath   string
	enableStrictKey  bool
	termHeight       *uint16
	termWidth        *uint16
}

func (o *transportBinOptions) apply( //nolint: gocyclo
	driverPtr uintptr,
	m *scrapligoffi.Mapping,
) error {
	if o.bin != "" {
		rc := m.Options.TransportBin.SetBin(driverPtr, o.bin)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting bin transport bin option", nil)
		}
	}

	if o.extraOpenArgs != "" {
		rc := m.Options.TransportBin.SetExtraOpenArgs(driverPtr, o.extraOpenArgs)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport extra args option",
				nil,
			)
		}
	}

	if o.overrideOpenArgs != "" {
		rc := m.Options.TransportBin.SetOverrideOpenArgs(driverPtr, o.overrideOpenArgs)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport override args option",
				nil,
			)
		}
	}

	if o.sshConfigPath != "" {
		rc := m.Options.TransportBin.SetSSHConfigPath(driverPtr, o.sshConfigPath)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport ssh config path option",
				nil,
			)
		}
	}

	if o.knownHostsPath != "" {
		rc := m.Options.TransportBin.SetKnownHostsPath(driverPtr, o.knownHostsPath)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport known hosts path option",
				nil,
			)
		}
	}

	if o.enableStrictKey {
		rc := m.Options.TransportBin.SetEnableStrictKey(driverPtr)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport strict key option",
				nil,
			)
		}
	}

	if o.termHeight != nil {
		rc := m.Options.TransportBin.SetTermHeight(driverPtr, *o.termHeight)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport term height option",
				nil,
			)
		}
	}

	if o.termWidth != nil {
		rc := m.Options.TransportBin.SetTermWidth(driverPtr, *o.termWidth)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport term width option",
				nil,
			)
		}
	}

	return nil
}

type transportSSH2Options struct {
	libSSH2Trace bool
}

func (o *transportSSH2Options) apply(driverPtr uintptr, m *scrapligoffi.Mapping) error {
	if o.libSSH2Trace {
		rc := m.Options.TransportSSH2.SetLibSSH2Trace(driverPtr)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting libssh2 trace option", nil)
		}
	}

	return nil
}

// WithLoggerCallback sets the logger callback for the Driver to use -- this is passed as pointer
// to the zig bits.
func WithLoggerCallback(loggerCallback func(level uint8, message *string)) Option {
	return func(d *Driver) error {
		d.options.loggerCallback = loggerCallback

		return nil
	}
}

// WithTransportKind sets the TransportKind to use in the Driver.
func WithTransportKind(transportKind TransportKind) Option {
	return func(d *Driver) error {
		d.options.transportKind = transportKind

		return nil
	}
}

// WithPort sets the port for the driver to connect to.
func WithPort(port uint16) Option {
	return func(d *Driver) error {
		d.options.port = &port

		return nil
	}
}
