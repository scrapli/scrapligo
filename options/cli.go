package options

import scrapligointernal "github.com/scrapli/scrapligo/internal"

// PlatformNameOrString is a string-like interface so you can pass a PlatformName or "normal" string
// to the driver constructor.
type PlatformNameOrString interface {
	~string
}

// WithSkipStaticOptions tells the Cli initialization to skip any "static" options if present --
// that is any options from scrapli_definitions that have been copied into scrapligo at release
// time and/or any options explicitly registered to the platform options singleton.
func WithSkipStaticOptions() Option {
	return func(o *scrapligointernal.Options) error {
		o.Cli.SkipStaticOptions = true

		return nil
	}
}

// WithDefintionFileOrName sets the Cli definition/platform for the Cli object.
func WithDefintionFileOrName[T PlatformNameOrString](s T) Option {
	return func(o *scrapligointernal.Options) error {
		o.Cli.DefinitionFileOrName = string(s)

		return nil
	}
}
