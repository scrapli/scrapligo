package cli

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// SendPromptedInput sends an `input` to the device expecting the given `prompt`, finally sending
// the `response`.
func (c *Cli) SendPromptedInput(
	ctx context.Context,
	input,
	prompt,
	response string,
	options ...OperationOption,
) (*Result, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	loadedOptions := newOperationOptions(options...)

	var operationID uint32

	status := c.ffiMap.Cli.SendPromptedInput(
		c.ptr,
		&operationID,
		&cancel,
		input,
		prompt,
		loadedOptions.promptPattern,
		response,
		loadedOptions.abortInput,
		loadedOptions.requestedMode,
		string(loadedOptions.inputHandling),
		loadedOptions.hiddenInput,
		loadedOptions.retainTrailingPrompt,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit sendPromptedInput operation", nil)
	}

	return c.getResult(ctx, &cancel, operationID)
}
