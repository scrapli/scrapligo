package options

import (
	"time"

	scrapligointernal "github.com/scrapli/scrapligo/internal"
	scrapligoutil "github.com/scrapli/scrapligo/util"
)

// WithReadSize sets the size of each individual read from the transport.
func WithReadSize(s uint64) Option {
	return func(o *scrapligointernal.Options) error {
		o.Session.ReadSize = &s

		return nil
	}
}

// WithReadDelayMin sets the minimum delay between reads -- this should be small but not crazy small
// otherwise cpu usage will suffer.
func WithReadDelayMin(t time.Duration) Option {
	return func(o *scrapligointernal.Options) error {
		v := scrapligoutil.SafeInt64ToUint64(t.Nanoseconds())
		o.Session.ReadDelayMinNs = &v

		return nil
	}
}

// WithReadDelayMax sets the maximum delay between reads. The minimum delay is backed off up to the
// maximum (this value) on subsequent reads that produce zero bytes. Once a read produces bytes
// again, the value is reset to the minimum.
func WithReadDelayMax(t time.Duration) Option {
	return func(o *scrapligointernal.Options) error {
		v := scrapligoutil.SafeInt64ToUint64(t.Nanoseconds())
		o.Session.ReadDelayMaxNs = &v

		return nil
	}
}

// WithReadDelayBackoffFactor sets the backoff factor from minimum to maximum read delay.
func WithReadDelayBackoffFactor(i uint8) Option {
	return func(o *scrapligointernal.Options) error {
		o.Session.ReadDelayBackoffFactor = &i

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
