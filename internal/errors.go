package internal

import (
	"context"
	"errors"
	"strings"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

// multiError wraps multiple errors such that errors.Is/errors.As match any of them, but unlike
// errors.Join, renders them on a single line.
type multiError struct {
	errs []error
}

func (e *multiError) Error() string {
	msgs := make([]string, len(e.errs))

	for idx, err := range e.errs {
		msgs[idx] = err.Error()
	}

	return strings.Join(msgs, ": ")
}

func (e *multiError) Unwrap() []error {
	return e.errs
}

// NewOperationError builds the error returned when libscrapli reports a failed operation. The
// errName is the libscrapli (zig) error name, and lastErr is the (possibly empty) human readable
// last error message from libscrapli. The returned error wraps the matching libscrapli sentinel
// error, and also wraps the context error if the context is done, so users can errors.Is against
// either scrapligoerrors.ErrTimeout or context.DeadlineExceeded (for example).
func NewOperationError(ctx context.Context, errName, lastErr string) error {
	sentinel := scrapligoerrors.LibscrapliErrorNameToError(errName)

	message := "operation failed"
	if errors.Is(sentinel, scrapligoerrors.ErrUnknown) {
		// include the raw name since it tells the user more than "unknown" does
		message += " (" + errName + ")"
	}

	if lastErr != "" {
		message += ": " + lastErr
	}

	var inner error = sentinel

	ctxErr := ctx.Err()
	if ctxErr != nil {
		inner = &multiError{errs: []error{sentinel, ctxErr}}
	}

	return scrapligoerrors.NewOperationError(message, inner)
}
