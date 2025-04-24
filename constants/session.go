package constants

const (
	// DefaultReadDelayMinNs is the default minimum read delay as set in the zig bits.
	DefaultReadDelayMinNs uint64 = 10_000
	// DefaultReadDelayMaxNs is the default maximum read delay as set in the zig bits.
	DefaultReadDelayMaxNs uint64 = 25_000_000
	// DefaultReadDelayBackoffFactor is the default read backoff factor as set in the zig bits.
	DefaultReadDelayBackoffFactor uint8 = 2
)
