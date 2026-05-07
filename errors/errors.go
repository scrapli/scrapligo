package errors

import (
	"errors"
	"fmt"
)

var (
	_ error        = (*ScrapliError)(nil)
	_ wrappedError = (*ScrapliError)(nil)
)

type wrappedError interface {
	Unwrap() error
}

// ErrNoMessages is an error returned when there are no more messages to check for a subscription or
// notification stream.
var ErrNoMessages = errors.New("no messages")

// ErrSubscriptionID is an error returned when failing to parse a subscription id from a message.
var ErrSubscriptionID = errors.New("subscription id")

// ErrorKind is an enum(ish) representing the kind of error -- i.e. "ffi" or "auth".
type ErrorKind string

const (
	// Ffi represents errors caused by the ffi integration layer, such as failure to submit or
	// poll an operation.
	Ffi ErrorKind = "ffi"
	// Options represents errors applying Cli options.
	Options ErrorKind = "options"
	// Netconf represents errors encountered during netconf operations.
	Netconf ErrorKind = "netconf"
	// Util represents errors encountered during utility funcs like parsing output.
	Util ErrorKind = "util"
)

// ScrapliError is the base error type used for all scrapli errors.
type ScrapliError struct {
	Kind    ErrorKind
	Message string
	Inner   error
}

func (e *ScrapliError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Inner)
	}

	return e.Message
}

func (e *ScrapliError) Unwrap() error {
	return e.Inner
}

func newScrapliError(kind ErrorKind, message string, inner error) error {
	return &ScrapliError{
		Kind:    kind,
		Message: message,
		Inner:   inner,
	}
}

// NewFfiError returns a "ffi" flavor ScrapliError, wrapping the inner error if provided.
func NewFfiError(message string, inner error) error {
	return newScrapliError(Ffi, message, inner)
}

// NewOptionsError returns an "options" flavor ScrapliError, wrapping the inner error if provided.
func NewOptionsError(message string, inner error) error {
	return newScrapliError(Options, message, inner)
}

// NewNetconfError returns a "netconf" flavor ScrapliError, wrapping the inner error if provided.
func NewNetconfError(message string, inner error) error {
	return newScrapliError(Netconf, message, inner)
}

// NewUtilError returns a "util" flavor ScrapliError, wrapping the inner error if provided.
func NewUtilError(message string, inner error) error {
	return newScrapliError(Util, message, inner)
}

// NewMessagesError returns a "netconf" flavor ScrapliError, wrapping the ErrNoMessages error type.
func NewMessagesError() error {
	return newScrapliError(Netconf, "no messages", ErrNoMessages)
}
