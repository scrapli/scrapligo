package driver

// WithUsername sets the username to use for authentication to the target device.
func WithUsername(username string) Option {
	return func(d *Driver) error {
		d.options.auth.username = username

		return nil
	}
}

// WithPassword sets the password to use for authentication to the target device.
func WithPassword(password string) Option {
	return func(d *Driver) error {
		d.options.auth.password = password

		return nil
	}
}

// WithAuthBypass bypasses/disables the "in session" authentication process where applicable (which
// means in the bin/telnet transports basically).
func WithAuthBypass() Option {
	return func(d *Driver) error {
		d.options.auth.inSessionAuthBypass = true

		return nil
	}
}

// WithUsernamePattern is a string that will be compiled via pcre2 in the underlying zig session
// object -- this pattern should match a username prompt for "in session" authentication (auth
// that happens "in" the session rather than in the transport natively (i.e. ssh2)).
func WithUsernamePattern(s string) Option {
	return func(d *Driver) error {
		d.options.auth.usernamePattern = s

		return nil
	}
}

// WithPasswordPattern is a string that will be compiled via pcre2 in the underlying zig session
// object -- this pattern should match a password prompt for "in session" authentication (auth
// that happens "in" the session rather than in the transport natively (i.e. ssh2)).
func WithPasswordPattern(s string) Option {
	return func(d *Driver) error {
		d.options.auth.passwordPattern = s

		return nil
	}
}

// WithPassphrasePattern is a string that will be compiled via pcre2 in the underlying zig session
// object -- this pattern should match a passphrase prompt for "in session" authentication (auth
// that happens "in" the session rather than in the transport natively (i.e. ssh2)).
func WithPassphrasePattern(s string) Option {
	return func(d *Driver) error {
		d.options.auth.passphrasePattern = s

		return nil
	}
}
