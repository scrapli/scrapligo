package options

import (
	"time"

	scrapligointernal "github.com/scrapli/scrapligo/internal"
	scrapligoutil "github.com/scrapli/scrapligo/util"
)

// WithReadSize sets the size of each individual read from the transport.
func WithReadSize(i uint64) Option {
	return func(o *scrapligointernal.Options) error {
		o.Session.ReadSize = &i

		return nil
	}
}

// WithReadMinDelay sets the minimum delay in ns between session reads.
func WithReadMinDelay(i uint64) Option {
	return func(o *scrapligointernal.Options) error {
		o.Session.ReadMinDelay = &i

		return nil
	}
}

// WithReadMaxDelay sets the maximum delay in ns between session reads.
func WithReadMaxDelay(i uint64) Option {
	return func(o *scrapligointernal.Options) error {
		o.Session.ReadMaxDelay = &i

		return nil
	}
}

// WithReturnChar sets the value to use as the return character -- normally this does not need to
// be set, however in some instances it may need to be set to carriage return + newline (\r\n)
// rather than the default newline (\n).
func WithReturnChar(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Session.ReturnChar = s

		return nil
	}
}

// WithOperationTimeout sets the default timeout for operations -- that is, unless otherwise
// specified on a given operation, this will be the timeout governing the operation.
func WithOperationTimeout(t time.Duration) Option {
	return func(o *scrapligointernal.Options) error {
		v := scrapligoutil.SafeInt64ToUint64(t.Nanoseconds())
		o.Session.OperationTimeoutNs = &v

		return nil
	}
}

// WithOperationMaxSearchDepth sets the maximum depth that a prompt is searched for when "looking"
// for the prompt pattern/delimiter -- you probably don't want to/need to change this. Making this
// unnecessarily large -- especially for "normal" ssh (not netconf) operations -- will slow things
// down quite a bit as that will be a larger blob to send to pcre2 for regex searching. Conversely,
// making this too small can lead to "missing" the prompt (because the pattern match may be
// incomplete).
func WithOperationMaxSearchDepth(i uint64) Option {
	return func(o *scrapligointernal.Options) error {
		o.Session.OperationMaxSearchDepth = &i

		return nil
	}
}

// WithSessionRecorderPath sets the output path for a recorder/writer for the session. DO NOT USE
// OTHER THAN FOR TESTING -- THIS IS UNSAFE AND WILL LEAK :).
func WithSessionRecorderPath(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Session.RecorderPath = s

		return nil
	}
}
