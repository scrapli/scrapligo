package driver

// WithBinTransportBinOverride overrides the default binary (/bin/ssh) used for the "bin" transport.
func WithBinTransportBinOverride(s string) Option {
	return func(d *Driver) error {
		d.options.transport.bin.bin = s

		return nil
	}
}

// WithBinTransportExtraArgs sets a list of additional arguments to add to the "open" command. This
// value is provided as a string the same way you would provide it when using ssh on the cli -- i.e.
// "-o ProxyCommand='foo' -P 1234" etc.
func WithBinTransportExtraArgs(s string) Option {
	return func(d *Driver) error {
		d.options.transport.bin.extraOpenArgs = s

		return nil
	}
}

// WithBinTransportOverrideArgs overrides open arguments with the provided options. This value is
// provided as a string the same way you would provide it when using ssh on the cli -- i.e.
// "-o ProxyCommand='foo' -P 1234" etc.
func WithBinTransportOverrideArgs(s string) Option {
	return func(d *Driver) error {
		d.options.transport.bin.overrideOpenArgs = s

		return nil
	}
}

// WithBinTransportSSHConfigFile sets the config file to use (via -F flag) for the command.
func WithBinTransportSSHConfigFile(s string) Option {
	return func(d *Driver) error {
		d.options.transport.bin.sshConfigPath = s

		return nil
	}
}

// WithBinTransportKnownHostsFile sets the known hosts file to use (via -o UserKnownHostsFile)
// for the command.
func WithBinTransportKnownHostsFile(s string) Option {
	return func(d *Driver) error {
		d.options.transport.bin.knownHostsPath = s

		return nil
	}
}

// WithBinTransportStrictKey enables strict key checking (via -o StrictHostKeyChecking=yes) for the
// command.
func WithBinTransportStrictKey() Option {
	return func(d *Driver) error {
		d.options.transport.bin.enableStrictKey = true

		return nil
	}
}

// WithTermHeight sets the size of terminal height for a given connection -- not applicable to all
// transports.
func WithTermHeight(i uint16) Option {
	return func(d *Driver) error {
		d.options.transport.bin.termHeight = &i

		return nil
	}
}

// WithTermWidth sets the size of terminal width for a given connection -- not applicable to all
// transports.
func WithTermWidth(i uint16) Option {
	return func(d *Driver) error {
		d.options.transport.bin.termWidth = &i

		return nil
	}
}
