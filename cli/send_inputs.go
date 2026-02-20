package cli

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligoutil "github.com/scrapli/scrapligo/util"
)

func newSendInputsOptions(options ...Option) *sendInputsOptions {
	o := &sendInputsOptions{
		inputHandling: InputHandlingFuzzy,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type sendInputsOptions struct {
	requestedMode          string
	inputHandling          InputHandling
	retainInput            bool
	retainTrailingPrompt   bool
	stopOnIndicatedFailure bool
}

// SendInputs send multiple "inputs" to the device.
func (c *Cli) SendInputs(
	ctx context.Context,
	inputs []string,
	options ...Option,
) (*Result, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	loadedOptions := newSendInputsOptions(options...)

	var results *Result

	for _, input := range inputs {
		var operationID uint32

		status := c.ffiMap.Cli.SendInput(
			c.ptr,
			&operationID,
			&cancel,
			input,
			loadedOptions.requestedMode,
			string(loadedOptions.inputHandling),
			loadedOptions.retainInput,
			loadedOptions.retainTrailingPrompt,
		)
		if status != 0 {
			return nil, scrapligoerrors.NewFfiError("failed to submit sendInput operation", nil)
		}

		r, err := c.getResult(ctx, &cancel, operationID)
		if err != nil {
			return nil, err
		}

		if results == nil {
			results = r
		} else {
			results.extend(r)
		}

		if r.Failed() && loadedOptions.stopOnIndicatedFailure {
			// note that this returns nil for an error since there was nothing unrecoverable
			// (probably) that happened, just we saw some stuff in the output saying that we
			// had a bad input or something
			return results, nil
		}
	}

	return results, nil
}

// SendInputsFromFile is a conveince wrapper to load inputs from a file then pass those to
// SendInputs.
func (c *Cli) SendInputsFromFile(
	ctx context.Context,
	f string,
	options ...Option,
) (*Result, error) {
	inputs, err := scrapligoutil.LoadFileLines(f)
	if err != nil {
		return nil, err
	}

	return c.SendInputs(ctx, inputs, options...)
}
