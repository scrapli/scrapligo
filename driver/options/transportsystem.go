package options

import (
	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util"
)

// WithSystemTransportOpenArgs sets the ExtraArgs value of the System Transport, these arguments are
// appended to the System transport open command (the command that spawns the connection).
func WithSystemTransportOpenArgs(l []string) util.Option {
	return func(o interface{}) error {
		t, ok := o.(*transport.System)

		if !ok {
			return util.ErrIgnoredOption
		}

		t.ExtraArgs = append(t.ExtraArgs, l...)

		return nil
	}
}
