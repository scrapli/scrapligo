package cli

import (
	"context"
	"strings"

	scrapligoconstants "github.com/scrapli/scrapligo/v2/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
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

	joinedInputs := strings.Join(inputs, scrapligoconstants.LibScrapliDelimiter)

	var operationID uint32

	err := c.ffiMap.Cli.SendInputs(
		c.ptr,
		&operationID,
		&cancel,
		joinedInputs,
		loadedOptions.requestedMode,
		string(loadedOptions.inputHandling),
		loadedOptions.retainInput,
		loadedOptions.retainTrailingPrompt,
	)
	if err != nil {
		return nil, err
	}

	return c.getResult(ctx, &cancel, operationID)
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
