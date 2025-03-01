package util

import "math"

// SafeInt64ToUint64 accepts an int64 and returns a uint64 -- if it would overflow (is negative) it
// simply returns the max uint64 value. So its lossy/can be bad, but won't break things.
func SafeInt64ToUint64(i int64) uint64 {
	if i < 0 {
		return 0
	}

	return uint64(i)
}

// SafeUint64ToInt64 accepts a uint64 and returns an int64 safely.
func SafeUint64ToInt64(i uint64) int64 {
	if i > math.MaxInt64 {
		return math.MaxInt64
	}

	return int64(i)
}
