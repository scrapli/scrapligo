package cli

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newSendPromptedInputOptions(options ...Option) *sendPromptedInputOptions {
	o := &sendPromptedInputOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type sendPromptedInputOptions struct {
	requestedMode        string
	inputHandling        *InputHandling
	retainTrailingPrompt bool
	promptPattern        string
	abortInput           string
	hiddenInput          bool
}

func (o *sendPromptedInputOptions) getInputHandling() *uint8 {
	if o.inputHandling == nil {
		return nil
	}

	v := uint8(*o.inputHandling)

	return &v
}

// SendPromptedInput sends an `input` to the device expecting the given `prompt`, finally sending
// the `response`.
func (c *Cli) SendPromptedInput(
	ctx context.Context,
	input,
	prompt,
	response string,
	options ...Option,
) (*Result, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	loadedOptions := newSendPromptedInputOptions(options...)

	var operationID uint32

	err := c.ffiMap.Cli.SendPromptedInput(
		c.ptr,
		&operationID,
		&cancel,
		input,
		prompt,
		loadedOptions.promptPattern,
		response,
		loadedOptions.abortInput,
		loadedOptions.requestedMode,
		loadedOptions.getInputHandling(),
		loadedOptions.hiddenInput,
		loadedOptions.retainTrailingPrompt,
	)
	if err != nil {
		return nil, err
	}

	return c.getResult(ctx, &cancel, operationID)
}
