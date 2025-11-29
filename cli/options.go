package cli

// InputHandling is an enum(ish) representing the kind of InputHandling an operation should use.
type InputHandling string

const (
	// InputHandlingExact represents "exact" input handling -- meaning the driver will read the
	// *exact* inputs from the connection before continuing/sending the return char.
	InputHandlingExact InputHandling = "exact"
	// InputHandlingFuzzy represents "fuzzy" input handling -- meaning the driver will read all
	// input characters in the correct order from the connection before continuing/sending return --
	// this means that inputs can still be read even if there are some chars that interrupt it such
	// as backspaces or indicators of newlines etc.
	InputHandlingFuzzy InputHandling = "fuzzy"
	// InputHandlingIgnore represents "ignore" input handling -- meaning the driver will simply
	// continue to sending return rather than attempting to read/consume the input. Generally,
	// don't use this.
	InputHandlingIgnore InputHandling = "ignore"
)

// Option defines a functional option for a cli operation.
type Option func(o any)

// WithRequestedMode sets the requested mode for the operation.
func WithRequestedMode(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case *sendInputOptions:
			to.requestedMode = s
		case *sendInputsOptions:
			to.requestedMode = s
		case *sendPromptedInputOptions:
			to.requestedMode = s
		}
	}
}

// WithInputHandling sets the input handling mode for the operation.
func WithInputHandling(i InputHandling) Option {
	return func(o any) {
		switch to := o.(type) {
		case *sendInputOptions:
			to.inputHandling = i
		case *sendInputsOptions:
			to.inputHandling = i
		case *sendPromptedInputOptions:
			to.inputHandling = i
		}
	}
}

// WithRetainInput retains (does not consume) the input of the operation. Ignored for
// SendPromptedInput.
func WithRetainInput() Option {
	return func(o any) {
		switch to := o.(type) {
		case *sendInputOptions:
			to.retainInput = true
		case *sendInputsOptions:
			to.retainInput = true
		}
	}
}

// WithRetainTrailingPrompt retains (does not consume) the trailing prompt at the end of the
// operations output.
func WithRetainTrailingPrompt() Option {
	return func(o any) {
		switch to := o.(type) {
		case *sendInputOptions:
			to.retainTrailingPrompt = true
		case *sendInputsOptions:
			to.retainTrailingPrompt = true
		case *sendPromptedInputOptions:
			to.retainTrailingPrompt = true
		}
	}
}

// WithStopOnIndicatedFailure stops sending inputs if some output indicates failure -- this is only
// applicable to the plural SendInputs operation, and is otherwise ignored.
func WithStopOnIndicatedFailure() Option {
	return func(o any) {
		switch to := o.(type) {
		case *sendInputsOptions:
			to.stopOnIndicatedFailure = true
		}
	}
}

// WithPromptPattern sets a string pcre2 regex pattern to look for after sending an input -- this is
// only applicable to SendPromptedInputs.
func WithPromptPattern(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case *sendPromptedInputOptions:
			to.promptPattern = s
		}
	}
}

// WithAbortInput sets a string to send to "abort" an operation if an indicated failure occurs --
// this is only applicable to SendPromptedInputs.
func WithAbortInput(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case *sendPromptedInputOptions:
			to.abortInput = s
		}
	}
}

// WithHiddenInput sets the "hidden input" option to true -- meaning the driver will expect to
// *not* be able to read the input being sent (likely because its a password prompt or something),
// this is only applicable to SendPromptedInputs.
func WithHiddenInput() Option {
	return func(o any) {
		switch to := o.(type) {
		case *sendPromptedInputOptions:
			to.hiddenInput = true
		}
	}
}

// WithContains sets the "contains" value of a read callback -- this is only applicable to
// ReadWithCallbacks callbacks (NewReadCallback).
func WithContains(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case *ReadCallback:
			to.contains = s
		}
	}
}

// WithContainsPattern sets the "contains pattern" value of a read callback -- this is only
// applicable to ReadWithCallbacks callbacks (NewReadCallback).
func WithContainsPattern(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case *ReadCallback:
			to.containsPattern = s
		}
	}
}

// WithNotContains sets the "not contains" value of a read callback -- this is only applicable to
// ReadWithCallbacks callbacks (NewReadCallback).
func WithNotContains(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case *ReadCallback:
			to.notContains = s
		}
	}
}

// WithSearchDepth sets the "search depth" value of a read callback -- this is only applicable to
// ReadWithCallbacks callbacks (NewReadCallback).
func WithSearchDepth(i uint64) Option {
	return func(o any) {
		switch to := o.(type) {
		case *ReadCallback:
			to.searchDepth = i
		}
	}
}

// WithOnce sets the "once" value of a read callback -- this is only applicable to
// ReadWithCallbacks callbacks (NewReadCallback).
func WithOnce() Option {
	return func(o any) {
		switch to := o.(type) {
		case *ReadCallback:
			to.once = true
		}
	}
}

// WithCompletes sets the "completes" value of a read callback -- this is only applicable to
// ReadWithCallbacks callbacks (NewReadCallback).
func WithCompletes() Option {
	return func(o any) {
		switch to := o.(type) {
		case *ReadCallback:
			to.completes = true
		}
	}
}
