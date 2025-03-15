package constants

const (
	// DefaultReadDelayMinNs is the default minimum read delay as set in the zig bits.
	DefaultReadDelayMinNs uint64 = 1_000
	// DefaultReadDelayMaxNs is the default maximum read delay as set in the zig bits.
	DefaultReadDelayMaxNs uint64 = 1_000_000
	// DefaultReadDelayBackoffFactor is the default read backoff factor as set in the zig bits.
	DefaultReadDelayBackoffFactor uint8 = 2

	// ReadDelayMultiplier is the multiplier that we apply to the read interval of the underlying
	// zig session object. By multiplying this we ensure that we dont have super tight loops in the
	// go bits and w aren't polling as fast as the zig bits can read. The zig read bits are suuuper
	// fast, so this being 2x is not a big deal.
	ReadDelayMultiplier = 2
)
