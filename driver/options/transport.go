package options

import (
	"time"

	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util"
)

// WithTransportReadSize sets the number of bytes each transport read operation should try to read.
// The default value is 65535.
func WithTransportReadSize(i int) util.Option {
	return func(o interface{}) error {
		a, ok := o.(*transport.Args)

		if !ok {
			return util.ErrIgnoredOption
		}

		a.ReadSize = i

		return nil
	}
}

// WithPort sets the TCP port for the connection, this defaults to 22 in all cases, so if using
// Telnet make sure you update the port!
func WithPort(i int) util.Option {
	return func(o interface{}) error {
		a, ok := o.(*transport.Args)

		if !ok {
			return util.ErrIgnoredOption
		}

		a.Port = i

		return nil
	}
}

// WithTermHeight sets the size to request for terminal (pty) height.
func WithTermHeight(i int) util.Option {
	return func(o interface{}) error {
		a, ok := o.(*transport.Args)

		if !ok {
			return util.ErrIgnoredOption
		}

		a.TermHeight = i

		return nil
	}
}

// WithTermWidth sets the size to request for terminal (pty) width.
func WithTermWidth(i int) util.Option {
	return func(o interface{}) error {
		a, ok := o.(*transport.Args)

		if !ok {
			return util.ErrIgnoredOption
		}

		a.TermWidth = i

		return nil
	}
}

// WithStandardTransportDialTimeout allows for modifying TimeoutSocket when using standard transport,
// this modifies the timeout for initial connections
func WithStandardTransportDialTimeout(t time.Duration) util.Option {
	return func(o interface{}) error {
		a, ok := o.(*transport.Args)
		if !ok {
			return util.ErrIgnoredOption
		}

		a.TimeoutSocket = t

		return nil
	}
}
