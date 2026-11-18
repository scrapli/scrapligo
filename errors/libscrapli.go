package errors

import "errors"

// ErrOutOfMemory is an error corresponding to the libscrapli result return code for out of memory
// errors.
var ErrOutOfMemory = errors.New("libscrapli: out of memory")

// ErrEOF is an error corresponding to the libscrapli result return code for eof.
var ErrEOF = errors.New("libscrapli: eof")

// ErrCancelled is an error corresponding to the libscrapli result return code for cancellation.
var ErrCancelled = errors.New("libscrapli: cancelled")

// ErrTimeout is an error corresponding to the libscrapli result return code for timeouts.
var ErrTimeout = errors.New("libscrapli: timeout")

// ErrDriver is an error corresponding to the libscrapli result return code for driver errors.
var ErrDriver = errors.New("libscrapli: driver")

// ErrSession is an error corresponding to the libscrapli result return code for session errors.
var ErrSession = errors.New("libscrapli: session")

// ErrTransport is an error corresponding to the libscrapli result return code for transport errors.
var ErrTransport = errors.New("libscrapli: transport")

// ErrOperation is an error corresponding to the libscrapli result return code for operation
// errors.
var ErrOperation = errors.New("libscrapli: operation")

// ErrInvalidArgument is an error corresponding to the libscrapli result return code for invalid
// arguments.
var ErrInvalidArgument = errors.New("libscrapli: invalid argument")

// ErrUnknown is an error corresponding to the libscrapli result return code for other errors
// that are not explicitly caught in libscrapli.
var ErrUnknown = errors.New("libscrapli: unknown/unhandled error")

// LibscrapliErrorNameToError returns the sentinel error for the given libscrapli error name (the
// zig error name, as returned by libscrapli for failed operations). This mirrors libscrapli's
// toFfiResult mapping (zig error -> FfiResult code) so that a given error resolves to the same
// sentinel whether it arrives as a return code or as an operation error name. Any name not
// handled there (including other ScrapliError members such as IndexError) returns ErrUnknown.
func LibscrapliErrorNameToError(name string) error {
	switch name {
	case "OutOfMemory":
		return ErrOutOfMemory
	case "EOF":
		return ErrEOF
	case "Cancelled":
		return ErrCancelled
	case "TimeoutExceeded":
		return ErrTimeout
	case "Driver":
		return ErrDriver
	case "Session":
		return ErrSession
	case "Transport":
		return ErrTransport
	case "Operation":
		return ErrOperation
	default:
		return ErrUnknown
	}
}
