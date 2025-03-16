package errors

var (
	_ error        = (*ScrapliError)(nil)
	_ wrappedError = (*ScrapliError)(nil)
)

type wrappedError interface {
	Unwrap() error
}

// ErrorKind is an enum(ish) representing the kind of error -- i.e. "ffi" or "auth".
type ErrorKind string

const (
	// Ffi represents errors caused by the ffi integration layer, such as failure to submit or
	// poll an operation.
	Ffi ErrorKind = "ffi"
	// Options represents errors applying Cli options.
	Options ErrorKind = "options"
	// Auth represents errors encountered during authentication.
	Auth ErrorKind = "auth"
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
	return e.Message
}

func (e *ScrapliError) Unwrap() error {
	return e.Inner
}

// NewFfiError returns a "ffi" flavor ScrapliError, wrapping the inner error if provided.
func NewFfiError(message string, inner error) error {
	e := &ScrapliError{
		Kind:    Ffi,
		Message: message,
	}

	if inner != nil {
		e.Inner = inner
	}

	return e
}

// NewOptionsError returns an "options" flavor ScrapliError, wrapping the inner error if provided.
func NewOptionsError(message string, inner error) error {
	e := &ScrapliError{
		Kind:    Options,
		Message: message,
	}

	if inner != nil {
		e.Inner = inner
	}

	return e
}

// NewAuthError returns an "auth" flavor ScrapliError, wrapping the inner error if provided.
func NewAuthError(message string, inner error) error {
	e := &ScrapliError{
		Kind:    Auth,
		Message: message,
	}

	if inner != nil {
		e.Inner = inner
	}

	return e
}

// NewUtilError returns a "util" flavor ScrapliError, wrapping the inner error if provided.
func NewUtilError(message string, inner error) error {
	e := &ScrapliError{
		Kind:    Util,
		Message: message,
	}

	if inner != nil {
		e.Inner = inner
	}

	return e
}
