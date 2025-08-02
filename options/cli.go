package options

import scrapligointernal "github.com/scrapli/scrapligo/internal"

// PlatformNameOrString is a string-like interface so you can pass a PlatformName or "normal" string
// to the driver constructor.
type PlatformNameOrString interface {
	~string
}

// WithDefintionFileOrName sets the Cli definition/platform for the Cli object.
func WithDefintionFileOrName[T PlatformNameOrString](s T) Option {
	return func(o *scrapligointernal.Options) error {
		o.Driver.DefinitionFileOrName = string(s)

		return nil
	}
}
