package options

import (
	"log"
	"log/slog"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
	scrapligointernal "github.com/scrapli/scrapligo/v2/internal"
	scrapligologging "github.com/scrapli/scrapligo/v2/logging"
)

// Option is a type used for functional options for the Cli object's options.
type Option func(o *scrapligointernal.Options) error

// WithLogger accepts one of:
//
// log.Logger
// slog.Logger
// func(level scrapligologging.LogLevel, message string)
//
//   - log.Logger will log everything with a level prefix (i.e. info, warn, etc.)
//   - slog.Logger will map scrapligologging.LogLevel to the most appropriate log function (i.e.
//     .Debug, .Info)
//   - the function will be invoked and passed the level/message if the cli/netconf objects
//     LoggingLevel is equal or higher
func WithLogger(
	logger any,
) Option {
	switch logger.(type) {
	case *log.Logger:
		return func(o *scrapligointernal.Options) error {
			o.Logger = logger

			return nil
		}
	case *slog.Logger:
		return func(o *scrapligointernal.Options) error {
			o.Logger = logger

			return nil
		}
	case func(scrapligologging.LogLevel, string):
		return func(o *scrapligointernal.Options) error {
			o.Logger = logger

			return nil
		}
	default:
		return func(_ *scrapligointernal.Options) error {
			return scrapligoerrors.NewOptionsError("invalid logger type provided", nil)
		}
	}
}

// WithLoggerLevel sets the log level for the given driver.
func WithLoggerLevel(
	level scrapligologging.LogLevel,
) Option {
	return func(o *scrapligointernal.Options) error {
		o.LoggerLevel = level

		return nil
	}
}

// WithPort sets the port for the driver to connect to.
func WithPort(port uint16) Option {
	return func(o *scrapligointernal.Options) error {
		o.Port = port

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
