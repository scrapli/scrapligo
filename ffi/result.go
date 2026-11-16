package ffi

import (
	"errors"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

const (
	LibscrapliReturnCodeSuccess uint8 = iota
	libscrapliReturnCodeUnknown
	libscrapliReturnCodeOutOfMemory
	libscrapliReturnCodeEOF
	libscrapliReturnCodeCancelled
	libscrapliReturnCodeTimeout
	libscrapliReturnCodeDriver
	libscrapliReturnCodeSession
	libscrapliReturnCodeTransport
	libscrapliReturnCodeOperation
	libscrapliReturnCodeInvalidArgument
)

// libscrapliResult holds a libscrapli return code and the caller of that functions provided
// error message and error factory -- the latter of which is used to return contextual information
// in the event the rc is a non-success value.
type libscrapliResult struct {
	rc      uint8
	message string
}

// newLibScrapliResult returns a LibscrapliResult wrapping a u8 return code from a libscrapli ffi
// call.
func newLibScrapliResult(
	rc uint8,
	message string,
) libscrapliResult {
	return libscrapliResult{
		rc:      rc,
		message: message,
	}
}

// check returns an error if the LibscrapliResult return code is a non-success value.
func (r libscrapliResult) check() error {
	if r.rc == LibscrapliReturnCodeSuccess {
		return nil
	}

	return scrapligoerrors.NewFfiError(r.message, resultCodeError(r.rc))
}

// resultCodeError returns the sentinel error matching a (non-success) libscrapli result code, or
// nil if the code does not map to a specific error.
func resultCodeError(rc uint8) error {
	switch rc {
	case libscrapliReturnCodeOutOfMemory:
		return scrapligoerrors.ErrOutOfMemory
	case libscrapliReturnCodeEOF:
		return scrapligoerrors.ErrEOF
	case libscrapliReturnCodeCancelled:
		return scrapligoerrors.ErrCancelled
	case libscrapliReturnCodeTimeout:
		return scrapligoerrors.ErrTimeout
	case libscrapliReturnCodeDriver:
		return scrapligoerrors.ErrDriver
	case libscrapliReturnCodeSession:
		return scrapligoerrors.ErrSession
	case libscrapliReturnCodeTransport:
		return scrapligoerrors.ErrTransport
	case libscrapliReturnCodeOperation:
		return scrapligoerrors.ErrOperation
	case libscrapliReturnCodeInvalidArgument:
		return scrapligoerrors.ErrInvalidArgument
	case libscrapliReturnCodeUnknown:
		return scrapligoerrors.ErrUnknown
	default:
		return nil
	}
}

// OperationError builds the error for a failed operation from the values fetched from libscrapli.
// The returned error wraps the sentinel error matching errCode (ErrOperation if the code does not
// map to a specific error) and ctxErr (if not nil) so users can errors.Is with either.
func OperationError(errCode uint8, errName, lastErr string, ctxErr error) error {
	message := errName
	if lastErr != "" {
		message += ": " + lastErr
	}

	inner := resultCodeError(errCode)
	if inner == nil || errors.Is(inner, scrapligoerrors.ErrUnknown) {
		inner = scrapligoerrors.ErrOperation
	}

	return scrapligoerrors.NewFfiError(message, errors.Join(inner, ctxErr))
}
