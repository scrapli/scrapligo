package internal

import (
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
)

// TransportKind is an enum(ish) representing the kind of transport a Cli should use.
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
	TransportKindTest TransportKind = "test_"
)

// NewOptions returns a new options object.
func NewOptions() *Options {
	return &Options{
		LoggerCallback: nil,
		TransportKind:  TransportKindBin,
		Port:           nil,
		Driver:         DriverOptions{},
		Netconf:        NetconfOptions{},
		Auth: AuthOptions{
			LookupMap: make(map[string]string),
		},
	}
}

// Options holds options for all driver kinds ("normal" and netconf).
type Options struct {
	Driver  DriverOptions
	Netconf NetconfOptions

	LoggerCallback func(level uint8, message *string)

	Port *uint16

	TransportKind TransportKind

	Session   SessionOptions
	Auth      AuthOptions
	Transport TransportOptions
}

// DriverOptions holds driver specific options.
type DriverOptions struct {
	DefinitionString string
}

// NetconfOptions holds netconf specific options.
type NetconfOptions struct{}

// Apply applies the Options to the given driver at driverPtr.
func (o *Options) Apply(driverPtr uintptr, m *scrapligoffi.Mapping) error {
	err := o.Session.apply(driverPtr, m)
	if err != nil {
		return err
	}

	err = o.Auth.apply(driverPtr, m)
	if err != nil {
		return err
	}

	switch o.TransportKind {
	case TransportKindBin:
		err = o.Transport.Bin.apply(driverPtr, m)
		if err != nil {
			return err
		}
	case TransportKindSSH2:
		err = o.Transport.SSH2.apply(driverPtr, m)
		if err != nil {
			return err
		}
	case TransportKindTelnet:
	case TransportKindTest:
		err = o.Transport.Test.apply(driverPtr, m)
		if err != nil {
			return err
		}
	}

	return nil
}

// SessionOptions holds options specific to the zig "Session" that lives in a driver.
type SessionOptions struct {
	ReadSize               *uint64
	ReadDelayMinNs         *uint64
	ReadDelayMaxNs         *uint64
	ReadDelayBackoffFactor *uint8
	ReturnChar             string

	OperationTimeoutNs      *uint64
	OperationMaxSearchDepth *uint64

	// do not use outside of tests, will leak/is unsafe!
	RecorderPath string
}

func (o *SessionOptions) apply( //nolint: gocyclo
	driverPtr uintptr,
	m *scrapligoffi.Mapping,
) error {
	if o.ReadSize != nil {
		rc := m.Options.Session.SetReadSize(driverPtr, *o.ReadSize)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting read size option", nil)
		}
	}

	if o.ReadDelayMinNs != nil {
		rc := m.Options.Session.SetReadDelayMinNs(driverPtr, *o.ReadDelayMinNs)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting min read delay option", nil)
		}
	}

	if o.ReadDelayMaxNs != nil {
		rc := m.Options.Session.SetReadDelayMaxNs(driverPtr, *o.ReadDelayMaxNs)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting max read delay option", nil)
		}
	}

	if o.ReadDelayBackoffFactor != nil {
		rc := m.Options.Session.SetReadDelayBackoffFactor(driverPtr, *o.ReadDelayBackoffFactor)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting read delay backoff factor option",
				nil,
			)
		}
	}

	if o.ReturnChar != "" {
		rc := m.Options.Session.SetReturnChar(driverPtr, o.ReturnChar)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting return char option", nil)
		}
	}

	if o.OperationTimeoutNs != nil {
		rc := m.Options.Session.SetOperationTimeoutNs(driverPtr, *o.OperationTimeoutNs)
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

	if o.OperationMaxSearchDepth != nil {
		rc := m.Options.Session.SetOperationMaxSearchDepth(driverPtr, *o.OperationMaxSearchDepth)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting operation search depth option",
				nil,
			)
		}
	}

	if o.RecorderPath != "" {
		rc := m.Options.Session.SetRecorderPath(driverPtr, o.RecorderPath)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting recorder path option", nil)
		}
	}

	return nil
}

// AuthOptions holds auth related options for driveres.
type AuthOptions struct {
	Username string
	Password string

	PrivateKeyPath       string
	PrivateKeyPassphrase string

	LookupMap map[string]string

	InSessionAuthBypass bool

	UsernamePattern   string
	PasswordPattern   string
	PassphrasePattern string
}

func (o *AuthOptions) apply(driverPtr uintptr, m *scrapligoffi.Mapping) error { //nolint: gocyclo
	if o.Username != "" {
		rc := m.Options.Auth.SetUsername(driverPtr, o.Username)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting username option", nil)
		}
	}

	if o.Password != "" {
		rc := m.Options.Auth.SetPassword(driverPtr, o.Password)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting password option", nil)
		}
	}

	if o.PrivateKeyPath != "" {
		rc := m.Options.Auth.SetPrivateKeyPath(driverPtr, o.PrivateKeyPath)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting private key path option", nil)
		}
	}

	if o.PrivateKeyPassphrase != "" {
		rc := m.Options.Auth.SetPrivateKeyPassphrase(driverPtr, o.PrivateKeyPassphrase)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting private key passphrase option",
				nil,
			)
		}
	}

	for k, v := range o.LookupMap {
		rc := m.Options.Auth.SetDriverOptionAuthLookupKeyValue(driverPtr, k, v)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting lookup map option",
				nil,
			)
		}
	}

	if o.InSessionAuthBypass {
		rc := m.Options.Auth.SetInSessionAuthBypass(driverPtr)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting in session auth bypass option",
				nil,
			)
		}
	}

	if o.UsernamePattern != "" {
		rc := m.Options.Auth.SetUsernamePattern(driverPtr, o.UsernamePattern)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting username pattern option", nil)
		}
	}

	if o.PasswordPattern != "" {
		rc := m.Options.Auth.SetPasswordPattern(driverPtr, o.PasswordPattern)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting password pattern option", nil)
		}
	}

	if o.PassphrasePattern != "" {
		rc := m.Options.Auth.SetPassphrasePattern(driverPtr, o.PassphrasePattern)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting passphrase pattern option", nil)
		}
	}

	return nil
}

// TransportOptions holds transport specific options.
type TransportOptions struct {
	Bin  TransportBinOptions
	SSH2 TransportSSH2Options
	Test TransportTestOptions
}

// TransportBinOptions holds "bin" transport specific options.
type TransportBinOptions struct {
	Bin              string
	ExtraOpenArgs    string
	OverrideOpenArgs string
	SSHConfigPath    string
	KnownHostsPath   string
	EnableStrictKey  bool
	TermHeight       *uint16
	TermWidth        *uint16
}

func (o *TransportBinOptions) apply( //nolint: gocyclo
	driverPtr uintptr,
	m *scrapligoffi.Mapping,
) error {
	if o.Bin != "" {
		rc := m.Options.TransportBin.SetBin(driverPtr, o.Bin)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting bin transport bin option", nil)
		}
	}

	if o.ExtraOpenArgs != "" {
		rc := m.Options.TransportBin.SetExtraOpenArgs(driverPtr, o.ExtraOpenArgs)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport extra args option",
				nil,
			)
		}
	}

	if o.OverrideOpenArgs != "" {
		rc := m.Options.TransportBin.SetOverrideOpenArgs(driverPtr, o.OverrideOpenArgs)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport override args option",
				nil,
			)
		}
	}

	if o.SSHConfigPath != "" {
		rc := m.Options.TransportBin.SetSSHConfigPath(driverPtr, o.SSHConfigPath)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport ssh config path option",
				nil,
			)
		}
	}

	if o.KnownHostsPath != "" {
		rc := m.Options.TransportBin.SetKnownHostsPath(driverPtr, o.KnownHostsPath)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport known hosts path option",
				nil,
			)
		}
	}

	if o.EnableStrictKey {
		rc := m.Options.TransportBin.SetEnableStrictKey(driverPtr)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport strict key option",
				nil,
			)
		}
	}

	if o.TermHeight != nil {
		rc := m.Options.TransportBin.SetTermHeight(driverPtr, *o.TermHeight)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport term height option",
				nil,
			)
		}
	}

	if o.TermWidth != nil {
		rc := m.Options.TransportBin.SetTermWidth(driverPtr, *o.TermWidth)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError(
				"failed setting bin transport term width option",
				nil,
			)
		}
	}

	return nil
}

// TransportSSH2Options holds (lib)"ssh2" transport specific options.
type TransportSSH2Options struct {
	LibSSH2Trace bool
}

func (o *TransportSSH2Options) apply(driverPtr uintptr, m *scrapligoffi.Mapping) error {
	if o.LibSSH2Trace {
		rc := m.Options.TransportSSH2.SetLibSSH2Trace(driverPtr)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting libssh2 trace option", nil)
		}
	}

	return nil
}

// TransportTestOptions holds test/file transport specific options.
type TransportTestOptions struct {
	F string
}

func (o *TransportTestOptions) apply(driverPtr uintptr, m *scrapligoffi.Mapping) error {
	if o.F != "" {
		rc := m.Options.TransportTest.SetF(driverPtr, o.F)
		if rc != 0 {
			return scrapligoerrors.NewOptionsError("failed setting test transport f option", nil)
		}
	}

	return nil
}
