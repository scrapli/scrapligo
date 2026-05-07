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
