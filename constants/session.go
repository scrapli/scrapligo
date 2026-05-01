package constants

const (
	// DefaultReadDelayMinNs is the default minimum read delay as set in the zig bits.
	DefaultReadDelayMinNs uint64 = 10_000
	// DefaultReadDelayMaxNs is the default maximum read delay as set in the zig bits.
	DefaultReadDelayMaxNs uint64 = 25_000_000

	// ReadyFDPollTimeoutNs is the timeout to use when polling the driver's ready signal.
	ReadyFDPollTimeoutNs int64 = 100_000_000
)
