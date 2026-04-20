package options

import scrapligointernal "github.com/kentik/scrapligo/v2/internal"

// WithTestTransportF sets the file to use for the test transport.
func WithTestTransportF(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Transport.Test.F = s

		return nil
	}
}
