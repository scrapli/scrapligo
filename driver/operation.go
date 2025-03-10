package driver

// OperationOption is a type used for functional options for Driver operations.
type OperationOption func(o *operationOptions)

// InputHandling is an enum(ish) representing the kind of InputHandling an operation should use.
type InputHandling string

const (
	// InputHandlingExact represents "exact" input handling -- meaning the driver will read the
	// *exact* inputs from the connection before continuing/sending the return char.
	InputHandlingExact InputHandling = "Exact"
	// InputHandlingFuzzy represents "fuzzy" input handling -- meaning the driver will read all
	// input characters in the correct order from the connection before continuing/sending return --
	// this means that inputs can still be read even if there are some chars that interrupt it such
	// as backspaces or indicators of newlines etc.
	InputHandlingFuzzy InputHandling = "Fuzzy"
	// InputHandlingIgnore represents "ignore" input handling -- meaning the driver will simply
	// continue to sending return rather than attempting to read/consume the input. Generally,
	// don't use this.
	InputHandlingIgnore InputHandling = "Ignore"
)

func newOperationOptions(options ...OperationOption) *operationOptions {
	o := &operationOptions{
		inputHandling: InputHandlingFuzzy,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type operationOptions struct {
	requestedMode          string
	inputHandling          InputHandling
	retainInput            bool
	retainTrailingPrompt   bool
	stopOnIndicatedFailure bool
	abortInput             string
	hiddenInput            bool
}

// WithRequestedMode sets the requested mode for the operation.
func WithRequestedMode(s string) OperationOption {
	return func(o *operationOptions) {
		o.requestedMode = s
	}
}

// WithInputHandling sets the input handling mode for the operation.
func WithInputHandling(i InputHandling) OperationOption {
	return func(o *operationOptions) {
		o.inputHandling = i
	}
}

// WithRetainInput retains (does not consume) the input of the operation. Ignored for
// SendPromptedInput.
func WithRetainInput() OperationOption {
	return func(o *operationOptions) {
		o.retainInput = true
	}
}

// WithRetainTrailingPrompt retains (does not consume) the trailing prompt at the end of the
// operations output.
func WithRetainTrailingPrompt() OperationOption {
	return func(o *operationOptions) {
		o.retainTrailingPrompt = true
	}
}

// WithStopOnIndicatedFailure stops sending inputs if some output indicates failure -- this is only
// applicable to the plural SendInputs operation, and is otherwise ignored.
func WithStopOnIndicatedFailure() OperationOption {
	return func(o *operationOptions) {
		o.stopOnIndicatedFailure = true
	}
}

// WithAbortInput sets a string to send to "abort" an operation if an indicated failure occurs --
// this is only applicable to SendPromptedInputs.
func WithAbortInput(s string) OperationOption {
	return func(o *operationOptions) {
		o.abortInput = s
	}
}

// WithHiddenInput sets the "hidden input" option to true -- meaning the driver will expect to
// *not* be able to read the input being sent (likely because its a password prompt or something),
// this is only applicable to SendPromptedInputs.
func WithHiddenInput() OperationOption {
	return func(o *operationOptions) {
		o.hiddenInput = true
	}
}
