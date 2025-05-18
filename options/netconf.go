package options

import scrapligointernal "github.com/scrapli/scrapligo/internal"

// NetconfVersion is an enumish type representing a netconf version (1.0 or 1.1).
type NetconfVersion string

const (
	// NetconfVersion10 represents the netconf version 1.0.
	NetconfVersion10 NetconfVersion = "1.0"
	// NetconfVersion11 represents the netconf version 1.1.
	NetconfVersion11 NetconfVersion = "1.1"
)

// WithNetconfErrorTag sets the error tag substring for a netconf object.
func WithNetconfErrorTag(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Netconf.ErrorTag = s

		return nil
	}
}

// WithNetconfPreferredVersion sets the preferred version for a netconf object.
func WithNetconfPreferredVersion(v NetconfVersion) Option {
	return func(o *scrapligointernal.Options) error {
		o.Netconf.PreferredVersion = string(v)

		return nil
	}
}

// WithNetconfMessagePollIntervalNS sets the message poll interval for a netconf object.
func WithNetconfMessagePollIntervalNS(v uint64) Option {
	return func(o *scrapligointernal.Options) error {
		o.Netconf.MessagePollIntervalNS = v

		return nil
	}
}
