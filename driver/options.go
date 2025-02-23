package driver

// Option is a type used for functional options for the Driver object's options.
type Option func(d *Driver) error

// TransportKind is an enum(ish) representing the kind of transport a Driver should use.
type TransportKind string

const (
	// TransportKindBin represents the "bin" transport -- the default transport that is a wrapper
	// around /bin/ssh.
	TransportKindBin TransportKind = "Bin"
	// TransportKindSSH2 represents the "ssh2" transport -- the transport using libssh2.
	TransportKindSSH2 TransportKind = "SSH2"
	// TransportKindTelnet represents the "telnet" transport.
	TransportKindTelnet TransportKind = "Telnet"
	// TransportKindFile represents the "file" transport that is used for integration testing.
	TransportKindFile TransportKind = "File"
)

const (
	// DefaultSSHPort is the default port used for SSH operations.
	DefaultSSHPort uint16 = 22
	// DefaultTelnetPort is the default port used for telnet operations.
	DefaultTelnetPort uint16 = 23
)

func newOptions() options {
	return options{
		platformVariant: "default",
		loggerCallback:  nil,
		transportKind:   TransportKindBin,
		port:            nil,
	}
}

type options struct {
	platformVariant string

	loggerCallback func(level uint8, message *string)

	transportKind TransportKind

	port     *uint16
	username string
	password string
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

// WithUsername sets the username to use for authentication to the target device.
func WithUsername(username string) Option {
	return func(d *Driver) error {
		d.options.username = username

		return nil
	}
}

// WithPassword sets the password to use for authentication to the target device.
func WithPassword(password string) Option {
	return func(d *Driver) error {
		d.options.password = password

		return nil
	}
}
