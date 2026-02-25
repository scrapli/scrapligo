package options

import scrapligointernal "github.com/scrapli/scrapligo/v2/internal"

// WithBinTransportBinOverride overrides the default binary (/bin/ssh) used for the "bin" transport.
func WithBinTransportBinOverride(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.Bin.Bin = s

		return nil
	}
}

// WithBinTransportExtraArgs sets a list of additional arguments to add to the "open" command. This
// value is provided as a string the same way you would provide it when using ssh on the cli -- i.e.
// "-o ProxyCommand='foo' -P 1234" etc.
func WithBinTransportExtraArgs(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.Bin.ExtraOpenArgs = s

		return nil
	}
}

// WithBinTransportOverrideArgs overrides open arguments with the provided options. This value is
// provided as a string the same way you would provide it when using ssh on the cli -- i.e.
// "-o ProxyCommand='foo' -P 1234" etc.
func WithBinTransportOverrideArgs(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.Bin.OverrideOpenArgs = s

		return nil
	}
}

// WithBinTransportSSHConfigFile sets the config file to use (via -F flag) for the command.
func WithBinTransportSSHConfigFile(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.Bin.SSHConfigPath = s

		return nil
	}
}

// WithBinTransportKnownHostsFile sets the known hosts file to use (via -o UserKnownHostsFile)
// for the command.
func WithBinTransportKnownHostsFile(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.Bin.KnownHostsPath = s

		return nil
	}
}

// WithBinTransportStrictKey enables strict key checking (via -o StrictHostKeyChecking=yes) for the
// command.
func WithBinTransportStrictKey() Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.Bin.EnableStrictKey = true

		return nil
	}
}

// WithTermHeight sets the size of terminal height for a given connection -- not applicable to all
// transports.
func WithTermHeight(i uint16) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.Bin.TermHeight = &i

		return nil
	}
}

// WithTermWidth sets the size of terminal width for a given connection -- not applicable to all
// transports.
func WithTermWidth(i uint16) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.Bin.TermWidth = &i

		return nil
	}
}
