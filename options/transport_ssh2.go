package options

import scrapligointernal "github.com/scrapli/scrapligo/v2/internal"

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

// WithSSH2ProxyJumpHost sets the end/final target host for a proxy jump style connection.
func WithSSH2ProxyJumpHost(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.SSH2.ProxyJumpHost = s

		return nil
	}
}

// WithSSH2ProxyJumpPort sets the end/final target port for a proxy jump style connection.
func WithSSH2ProxyJumpPort(i uint16) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.SSH2.ProxyJumpPort = i

		return nil
	}
}

// WithSSH2ProxyJumpUsername sets the end/final target username for a proxy jump style connection.
func WithSSH2ProxyJumpUsername(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.SSH2.ProxyJumpUsername = s

		return nil
	}
}

// WithSSH2ProxyJumpPassword sets the end/final target password for a proxy jump style connection.
func WithSSH2ProxyJumpPassword(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.SSH2.ProxyJumpPassword = s

		return nil
	}
}

// WithSSH2ProxyJumpPrivateKeyPath sets the end/final target private key path for a proxy jump
// style connection.
func WithSSH2ProxyJumpPrivateKeyPath(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.SSH2.ProxyJumpPrivateKeyPassphrase = s

		return nil
	}
}

// WithSSH2ProxyJumpPrivateKeyPassphrase sets the end/final target private key passphrase for a
// proxy jump style connection.
func WithSSH2ProxyJumpPrivateKeyPassphrase(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.SSH2.ProxyJumpPrivateKeyPassphrase = s

		return nil
	}
}

// WithSSH2ProxyJumpLibssh2Trace enables the libssh2 trace setup for a proxy jump style connection.
func WithSSH2ProxyJumpLibssh2Trace() Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.SSH2.ProxyJumpLibSSH2Trace = true

		return nil
	}
}
