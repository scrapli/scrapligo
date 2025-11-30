package internal

import (
	"unsafe"

	"github.com/ebitengine/purego"
	scrapligologging "github.com/scrapli/scrapligo/logging"
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

var zero uint64 //nolint: gochecknoglobals

// Options holds options for all driver kinds (cli and netconf).
type Options struct {
	Cli     CliOptions
	Netconf NetconfOptions

	Logger      any
	LoggerLevel scrapligologging.LogLevel

	Port uint16

	TransportKind TransportKind

	Session   SessionOptions
	Auth      AuthOptions
	Transport TransportOptions
}

// NewOptions returns a new options object.
func NewOptions() *Options {
	return &Options{
		Logger:        nil,
		LoggerLevel:   scrapligologging.Warn,
		TransportKind: TransportKindBin,
		Port:          0,
		Cli: CliOptions{
			DefinitionFileOrName: "default",
		},
		Netconf: NetconfOptions{},
		Auth: AuthOptions{
			LookupMap: make(map[string]string),
		},
	}
}

// GetLogger returns the AnyLogger wrapped around the configured logger options.
func (o *Options) GetLogger() *scrapligologging.AnyLogger {
	return scrapligologging.LoggerToAnyLogger(o.Logger, o.LoggerLevel)
}

// Apply applies the Options to the given driver options struct at optionsPtr.
func (o *Options) Apply(optionsPtr uintptr) {
	opts := (*driverOptions)(unsafe.Pointer(optionsPtr)) //nolint: govet

	opts.loggerLevel = uintptr(unsafe.Pointer(&[]byte(o.LoggerLevel)[0]))
	opts.loggerLevelLen = uintptr(len(o.LoggerLevel))
	opts.loggerCallback = scrapligologging.LoggerToLoggerCallback(
		o.Logger,
		uint8(scrapligologging.IntFromLevel(o.LoggerLevel)),
	)

	opts.port = &o.Port

	o.Cli.apply(opts)
	o.Netconf.apply(opts)
	o.Session.apply(opts)
	o.Auth.apply(opts)

	opts.transportKind = uintptr(unsafe.Pointer(&[]byte(o.TransportKind)[0]))
	opts.transportKindLen = uintptr(len(o.TransportKind))

	switch o.TransportKind {
	case TransportKindBin:
		o.Transport.Bin.apply(opts)
	case TransportKindSSH2:
		o.Transport.SSH2.apply(opts)
	case TransportKindTelnet:
	case TransportKindTest:
		o.Transport.Test.apply(opts)
	}
}

// CliOptions holds cli specific options.
type CliOptions struct {
	DefinitionFileOrName string
	DefinitionString     string
	SkipStaticOptions    bool
}

func (o *CliOptions) apply(opts *driverOptions) {
	if o.DefinitionString == "" {
		return
	}

	opts.cli.definitionStr = uintptr(unsafe.Pointer(&[]byte(o.DefinitionString)[0]))
	opts.cli.definitionStrLen = uintptr(len(o.DefinitionString))
}

// NetconfOptions holds netconf specific options.
type NetconfOptions struct {
	ErrorTag              string
	PreferredVersion      string
	MessagePollIntervalNS uint64
}

func (o *NetconfOptions) apply(opts *driverOptions) {
	if o.ErrorTag != "" {
		opts.netconf.errorTag = uintptr(unsafe.Pointer(&[]byte(o.ErrorTag)[0]))
		opts.netconf.errorTagLen = uintptr(len(o.ErrorTag))
	}

	if o.PreferredVersion != "" {
		opts.netconf.preferredVersion = uintptr(unsafe.Pointer(&[]byte(o.PreferredVersion)[0]))
		opts.netconf.preferredVersionLen = uintptr(len(o.PreferredVersion))
	}

	if o.MessagePollIntervalNS != 0 {
		opts.netconf.messagePollInterval = &o.MessagePollIntervalNS
	}
}

// SessionOptions holds options specific to the zig "Session" that lives in a driver.
type SessionOptions struct {
	ReadSize     *uint64
	ReadMinDelay *uint64
	ReadMaxDelay *uint64
	ReturnChar   string

	OperationTimeoutNs      *uint64
	OperationMaxSearchDepth *uint64

	RecorderPath     string
	RecorderCallback func(buf *[]byte)
}

func (o *SessionOptions) apply(opts *driverOptions) {
	if o.ReadSize != nil {
		opts.session.readSize = o.ReadSize
	}

	if o.ReadMinDelay != nil {
		opts.session.readMinDelayNs = o.ReadMinDelay
	}

	if o.ReadMaxDelay != nil {
		opts.session.readMaxDelayNs = o.ReadMaxDelay
	}

	if o.ReturnChar != "" {
		opts.session.returnChar = uintptr(unsafe.Pointer(&[]byte(o.ReturnChar)[0]))
		opts.session.returnCharLen = uintptr(len(o.ReturnChar))
	}

	if o.OperationTimeoutNs != nil {
		opts.session.operationTimeoutNs = o.OperationTimeoutNs
	} else {
		// if user does not provide a timeout we assume they want to govern all timeouts via context
		// cancellation
		opts.session.operationTimeoutNs = &zero
	}

	if o.OperationMaxSearchDepth != nil {
		opts.session.operationMaxSearchDepth = o.OperationMaxSearchDepth
	}

	if o.RecorderPath != "" {
		opts.session.recordDestination = uintptr(unsafe.Pointer(&[]byte(o.RecorderPath)[0]))
		opts.session.recordDestinationLen = uintptr(len(o.RecorderPath))
	} else if o.RecorderCallback != nil {
		opts.session.recorderCallback = purego.NewCallback(o.RecorderCallback)
	}
}

// AuthOptions holds auth related options for driveres.
type AuthOptions struct {
	Username string
	Password string

	PrivateKeyPath       string
	PrivateKeyPassphrase string

	LookupMap        map[string]string
	lookupMapKeys    []string
	lookupMapKeyLens []uint16
	lookupMapVals    []string
	lookupMapValLens []uint16

	ForceInSessionAuth  bool
	BypassInSessionAuth bool

	UsernamePattern   string
	PasswordPattern   string
	PassphrasePattern string
}

func (o *AuthOptions) apply(opts *driverOptions) {
	if o.Username != "" {
		opts.auth.username = uintptr(unsafe.Pointer(&[]byte(o.Username)[0]))
		opts.auth.usernameLen = uintptr(len(o.Username))
	}

	if o.Password != "" {
		opts.auth.password = uintptr(unsafe.Pointer(&[]byte(o.Password)[0]))
		opts.auth.passwordLen = uintptr(len(o.Password))
	}

	if o.PrivateKeyPath != "" {
		opts.auth.privateKeyPath = uintptr(unsafe.Pointer(&[]byte(o.PrivateKeyPath)[0]))
		opts.auth.privateKeyPathLen = uintptr(len(o.PrivateKeyPath))
	}

	if o.PrivateKeyPassphrase != "" {
		opts.auth.privateKeyPassphrase = uintptr(unsafe.Pointer(&[]byte(o.PrivateKeyPassphrase)[0]))
		opts.auth.privateKeyPassphraseLen = uintptr(len(o.PrivateKeyPassphrase))
	}

	if len(o.LookupMap) > 0 {
		// this ensures that the string/uint16 slices have a lifetime that lasts as long as
		// `o` which we know will last as long as the option apply process in zig
		o.lookupMapKeys = make([]string, len(o.LookupMap))
		o.lookupMapKeyLens = make([]uint16, len(o.LookupMap))
		o.lookupMapVals = make([]string, len(o.LookupMap))
		o.lookupMapValLens = make([]uint16, len(o.LookupMap))

		var count uint16

		for k, v := range o.LookupMap {
			o.lookupMapKeys[count] = k
			o.lookupMapKeyLens[count] = uint16(len(k)) //nolint: gosec

			o.lookupMapVals[count] = v
			o.lookupMapValLens[count] = uint16(len(v)) //nolint: gosec

			count++
		}

		opts.auth.lookups.keys = uintptr(unsafe.Pointer(&o.lookupMapKeys[0]))
		opts.auth.lookups.keysLens = uintptr(unsafe.Pointer(&o.lookupMapKeyLens[0]))

		opts.auth.lookups.vals = uintptr(unsafe.Pointer(&o.lookupMapVals[0]))
		opts.auth.lookups.valsLens = uintptr(unsafe.Pointer(&o.lookupMapValLens[0]))

		opts.auth.lookups.count = count
	}

	if o.ForceInSessionAuth {
		opts.auth.forceInSessionAuth = &o.ForceInSessionAuth
	}

	if o.BypassInSessionAuth {
		opts.auth.bypassInSessionAuth = &o.BypassInSessionAuth
	}

	if o.UsernamePattern != "" {
		opts.auth.usernamePattern = uintptr(unsafe.Pointer(&[]byte(o.UsernamePattern)[0]))
		opts.auth.usernamePatternLen = uintptr(len(o.UsernamePattern))
	}

	if o.PasswordPattern != "" {
		opts.auth.passwordPattern = uintptr(unsafe.Pointer(&[]byte(o.PasswordPattern)[0]))
		opts.auth.passwordPatternLen = uintptr(len(o.PasswordPattern))
	}

	if o.PassphrasePattern != "" {
		opts.auth.privateKeyPassphrasePattern = uintptr(
			unsafe.Pointer(&[]byte(o.PassphrasePattern)[0]),
		)
		opts.auth.privateKeyPassphrasePatternLen = uintptr(len(o.PassphrasePattern))
	}
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

func (o *TransportBinOptions) apply(opts *driverOptions) {
	if o.Bin != "" {
		opts.transport.bin.bin = uintptr(unsafe.Pointer(&[]byte(o.Bin)[0]))
		opts.transport.bin.binLen = uintptr(len(o.Bin))
	}

	if o.ExtraOpenArgs != "" {
		opts.transport.bin.extraOpenArgs = uintptr(unsafe.Pointer(&[]byte(o.ExtraOpenArgs)[0]))
		opts.transport.bin.extraOpenArgsLen = uintptr(len(o.ExtraOpenArgs))
	}

	if o.OverrideOpenArgs != "" {
		opts.transport.bin.overrideOpenArgs = uintptr(
			unsafe.Pointer(&[]byte(o.OverrideOpenArgs)[0]),
		)
		opts.transport.bin.overrideOpenArgsLen = uintptr(len(o.OverrideOpenArgs))
	}

	if o.SSHConfigPath != "" {
		opts.transport.bin.sshConfigPath = uintptr(unsafe.Pointer(&[]byte(o.SSHConfigPath)[0]))
		opts.transport.bin.sshConfigPathLen = uintptr(len(o.SSHConfigPath))
	}

	if o.KnownHostsPath != "" {
		opts.transport.bin.knownHostsPath = uintptr(unsafe.Pointer(&[]byte(o.KnownHostsPath)[0]))
		opts.transport.bin.knownHostsPathLen = uintptr(len(o.KnownHostsPath))
	}

	if o.EnableStrictKey {
		opts.transport.bin.enableStrictKey = &o.EnableStrictKey
	}

	if o.TermHeight != nil {
		opts.transport.bin.termHeight = o.TermHeight
	}

	if o.TermWidth != nil {
		opts.transport.bin.termWidth = o.TermWidth
	}
}

// TransportSSH2Options holds (lib)"ssh2" transport specific options.
type TransportSSH2Options struct {
	KnownHostsPath string
	LibSSH2Trace   bool

	// proxy jump related
	ProxyJumpHost                 string
	ProxyJumpPort                 uint16
	ProxyJumpUsername             string
	ProxyJumpPassword             string
	ProxyJumpPrivateKeyPath       string
	ProxyJumpPrivateKeyPassphrase string
	ProxyJumpLibSSH2Trace         bool
}

func (o *TransportSSH2Options) apply(opts *driverOptions) {
	if o.KnownHostsPath != "" {
		opts.transport.ssh2.knownHostsPath = uintptr(unsafe.Pointer(&[]byte(o.KnownHostsPath)[0]))
		opts.transport.ssh2.knownHostsPathLen = uintptr(len(o.KnownHostsPath))
	}

	if o.LibSSH2Trace {
		opts.transport.ssh2.libssh2Trace = &o.LibSSH2Trace
	}

	if o.ProxyJumpHost != "" {
		opts.transport.ssh2.proxyJumpHost = uintptr(unsafe.Pointer(&[]byte(o.ProxyJumpHost)[0]))
		opts.transport.ssh2.proxyJumpHostLen = uintptr(len(o.ProxyJumpHost))
	}

	if o.ProxyJumpPort != 0 {
		opts.transport.ssh2.proxyJumpPort = &o.ProxyJumpPort
	}

	if o.ProxyJumpUsername != "" {
		opts.transport.ssh2.proxyJumpUsername = uintptr(
			unsafe.Pointer(&[]byte(o.ProxyJumpUsername)[0]),
		)
		opts.transport.ssh2.proxyJumpUsernameLen = uintptr(len(o.ProxyJumpUsername))
	}

	if o.ProxyJumpPassword != "" {
		opts.transport.ssh2.proxyJumpPassword = uintptr(
			unsafe.Pointer(&[]byte(o.ProxyJumpPassword)[0]),
		)
		opts.transport.ssh2.proxyJumpPasswordLen = uintptr(len(o.ProxyJumpPassword))
	}

	if o.ProxyJumpPrivateKeyPath != "" {
		opts.transport.ssh2.proxyJumpPrivateKeyPath = uintptr(
			unsafe.Pointer(&[]byte(o.ProxyJumpPrivateKeyPath)[0]),
		)
		opts.transport.ssh2.proxyJumpPrivateKeyPathLen = uintptr(len(o.ProxyJumpPrivateKeyPath))
	}

	if o.ProxyJumpPrivateKeyPassphrase != "" {
		opts.transport.ssh2.proxyJumpPrivateKeyPassphrase = uintptr(
			unsafe.Pointer(&[]byte(o.ProxyJumpPrivateKeyPassphrase)[0]),
		)
		opts.transport.ssh2.proxyJumpPrivateKeyPassphraseLen = uintptr(
			len(o.ProxyJumpPrivateKeyPassphrase),
		)
	}

	if o.ProxyJumpLibSSH2Trace {
		opts.transport.ssh2.proxyJumpLibssh2Trace = &o.ProxyJumpLibSSH2Trace
	}
}

// TransportTestOptions holds test/file transport specific options.
type TransportTestOptions struct {
	F string
}

func (o *TransportTestOptions) apply(opts *driverOptions) {
	if o.F != "" {
		opts.transport.test.f = uintptr(unsafe.Pointer(&[]byte(o.F)[0]))
		opts.transport.test.fLen = uintptr(len(o.F))
	}
}
