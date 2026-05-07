package ffi

import (
	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

const (
	libscrapliReturnCodeSuccess uint8 = iota
	libscrapliReturnCodeOutOfMemory
	libscrapliReturnCodeEOF
	libscrapliReturnCodeCancelled
	libscrapliReturnCodeTimeout
	libscrapliReturnCodeDriver
	libscrapliReturnCodeSession
	libscrapliReturnCodeTransport
	libscrapliReturnCodeOperation
	libscrapliReturnCodeInvalidArgument
	libscrapliReturnCodeUnknown
)

// libscrapliResult holds a libscrapli return code and the caller of that functions provided
// error message and error factory -- the latter of which is used to return contextual information
// in the event the rc is a non-success value.
type libscrapliResult struct {
	rc                uint8
	message           string
	defaultErrFactory func(message string, inner error) error
}

// newLibScrapliResult returns a LibscrapliResult wrapping a u8 return code from a libscrapli ffi
// call.
func newLibScrapliResult(
	rc uint8,
	message string,
	defaultErrFactory func(message string, inner error) error,
) libscrapliResult {
	return libscrapliResult{
		rc:                rc,
		message:           message,
		defaultErrFactory: defaultErrFactory,
	}
}

// check returns an error if the LibscrapliResult return code is a non-success value.
func (r libscrapliResult) check() error {
	if r.rc == libscrapliReturnCodeSuccess {
		return nil
	}

	var inner error

	switch r.rc {
	case libscrapliReturnCodeOutOfMemory:
		inner = scrapligoerrors.ErrOutOfMemory
	case libscrapliReturnCodeEOF:
		inner = scrapligoerrors.ErrEOF
	case libscrapliReturnCodeCancelled:
		inner = scrapligoerrors.ErrCancelled
	case libscrapliReturnCodeTimeout:
		inner = scrapligoerrors.ErrTimeout
	case libscrapliReturnCodeDriver:
		inner = scrapligoerrors.ErrDriver
	case libscrapliReturnCodeSession:
		inner = scrapligoerrors.ErrSession
	case libscrapliReturnCodeTransport:
		inner = scrapligoerrors.ErrTransport
	case libscrapliReturnCodeOperation:
		inner = scrapligoerrors.ErrOperation
	case libscrapliReturnCodeInvalidArgument:
		inner = scrapligoerrors.ErrInvalidArgument
	case libscrapliReturnCodeUnknown:
		inner = scrapligoerrors.ErrUnknown
	}

	if r.defaultErrFactory != nil {
		return r.defaultErrFactory(r.message, inner)
	}

	return scrapligoerrors.NewFfiError(r.message, inner)
}
