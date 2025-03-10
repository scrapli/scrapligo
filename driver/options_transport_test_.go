package driver

// WithTestTransportF sets the file to use for the test transport.
func WithTestTransportF(s string) Option {
	return func(d *Driver) error {
		d.options.transport.test.f = s

		return nil
	}
}
