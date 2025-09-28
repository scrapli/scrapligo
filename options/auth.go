package options

import scrapligointernal "github.com/scrapli/scrapligo/internal"

// WithUsername sets the username to use for authentication to the target device.
func WithUsername(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Auth.Username = s

		return nil
	}
}

// WithPassword sets the password to use for authentication to the target device.
func WithPassword(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Auth.Password = s

		return nil
	}
}

// WithPrivateKeyPath sets the private key to use for authentication to the target device.
func WithPrivateKeyPath(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Auth.PrivateKeyPath = s

		return nil
	}
}

// WithPrivateKeyPassphrase sets the private key passhrase to use for authentication to the target
// device.
func WithPrivateKeyPassphrase(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Auth.PrivateKeyPassphrase = s

		return nil
	}
}

// WithLookupKeyValue adds an entry to the lookup map for the driver instance.
func WithLookupKeyValue(key, value string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Auth.LookupMap[key] = value

		return nil
	}
}

// WithForceInSessionAuth unconditionally forces the in session auth to run.
func WithForceInSessionAuth() Option {
	return func(o *scrapligointernal.Options) error {
		o.Auth.ForceInSessionAuth = true

		return nil
	}
}

// WithBypassInSessionAuth bypasses/disables the "in session" authentication process where
// applicable (which means in the bin/telnet transports basically).
func WithBypassInSessionAuth() Option {
	return func(o *scrapligointernal.Options) error {
		o.Auth.BypassInSessionAuth = true

		return nil
	}
}

// WithUsernamePattern is a string that will be compiled via pcre2 in the underlying zig session
// object -- this pattern should match a username prompt for "in session" authentication (auth
// that happens "in" the session rather than in the transport natively (i.e. ssh2)).
func WithUsernamePattern(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Auth.UsernamePattern = s

		return nil
	}
}

// WithPasswordPattern is a string that will be compiled via pcre2 in the underlying zig session
// object -- this pattern should match a password prompt for "in session" authentication (auth
// that happens "in" the session rather than in the transport natively (i.e. ssh2)).
func WithPasswordPattern(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Auth.PasswordPattern = s

		return nil
	}
}

// WithPassphrasePattern is a string that will be compiled via pcre2 in the underlying zig session
// object -- this pattern should match a passphrase prompt for "in session" authentication (auth
// that happens "in" the session rather than in the transport natively (i.e. ssh2)).
func WithPassphrasePattern(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Auth.PassphrasePattern = s

		return nil
	}
}
