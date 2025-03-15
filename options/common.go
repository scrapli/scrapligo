package options

import scrapligointernal "github.com/scrapli/scrapligo/internal"

// Option is a type used for functional options for the Driver object's options.
type Option func(o *scrapligointernal.Options) error

// WithLoggerCallback sets the logger callback for the Driver to use -- this is passed as pointer
// to the zig bits.
func WithLoggerCallback(
	loggerCallback func(level uint8, message *string),
) Option {
	return func(o *scrapligointernal.Options) error {
		o.LoggerCallback = loggerCallback

		return nil
	}
}

// WithPort sets the port for the driver to connect to.
func WithPort(port uint16) Option {
	return func(o *scrapligointernal.Options) error {
		o.Port = &port

		return nil
	}
}

// WithTransportBin sets the transport kind to "bin".
func WithTransportBin() Option {
	return func(o *scrapligointernal.Options) error {
		o.TransportKind = scrapligointernal.TransportKindBin

		return nil
	}
}

// WithTransportSSH2 sets the transport kind to "ssh2".
func WithTransportSSH2() Option {
	return func(o *scrapligointernal.Options) error {
		o.TransportKind = scrapligointernal.TransportKindSSH2

		return nil
	}
}

// WithTransportTelnet sets the transport kind to "telnet".
func WithTransportTelnet() Option {
	return func(o *scrapligointernal.Options) error {
		o.TransportKind = scrapligointernal.TransportKindTelnet

		return nil
	}
}

// WithTransportTest sets the transport kind to "test_".
func WithTransportTest() Option {
	return func(o *scrapligointernal.Options) error {
		o.TransportKind = scrapligointernal.TransportKindTest

		return nil
	}
}
