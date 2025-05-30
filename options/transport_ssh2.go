package options

import scrapligointernal "github.com/scrapli/scrapligo/internal"

// WithSSH2KnownHostsPath sets the known hosts file to use with libssh2 connections.
func WithSSH2KnownHostsPath(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.SSH2.KnownHostsPath = s

		return nil
	}
}

// WithSSH2LibSSH2Trace enables libssh2 trace logging.
func WithSSH2LibSSH2Trace() Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.SSH2.LibSSH2Trace = true

		return nil
	}
}
