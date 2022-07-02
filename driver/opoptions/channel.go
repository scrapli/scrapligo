package opoptions

import (
	"regexp"
	"time"

	"github.com/scrapli/scrapligo/driver/netconf"

	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/util"
)

// WithNoStripPrompt disables stripping the prompt out from the read bytes.
func WithNoStripPrompt() util.Option {
	return func(o interface{}) error {
		c, ok := o.(*channel.OperationOptions)

		if ok {
			c.StripPrompt = false

			return nil
		}

		return util.ErrIgnoredOption
	}
}

// WithEager forces the channel read operation into "eager" mode -- that is, it will no longer read
// inputs off of the channel prior to sending a return, hence "eager".
func WithEager() util.Option {
	return func(o interface{}) error {
		c, ok := o.(*channel.OperationOptions)

		if ok {
			c.Eager = true

			return nil
		}

		return util.ErrIgnoredOption
	}
}

// WithTimeoutOps modifies the timeout "ops" value, or the timeout for a given operation. This only
// modifies the timeout for the current operation and does not update the actual Channel TimeoutOps
// value permanently.
func WithTimeoutOps(t time.Duration) util.Option {
	return func(o interface{}) error {
		switch oo := o.(type) {
		case *channel.OperationOptions:
			oo.Timeout = t
		case *netconf.OperationOptions:
			oo.Timeout = t
		default:
			return util.ErrIgnoredOption
		}

		return nil
	}
}

// WithCompletePatterns is a slice of regex patterns that, if seen, indicate that the operation is
// complete -- this is used with SendInteractive.
func WithCompletePatterns(p []*regexp.Regexp) util.Option {
	return func(o interface{}) error {
		c, ok := o.(*channel.OperationOptions)

		if ok {
			c.CompletePatterns = p

			return nil
		}

		return util.ErrIgnoredOption
	}
}
