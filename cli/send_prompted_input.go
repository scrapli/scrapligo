package cli

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// SendPromptedInput sends an `input` to the device expecting the given `prompt`, finally sending
// the `response`.
func (d *Driver) SendPromptedInput(
	ctx context.Context,
	input,
	prompt,
	response string,
	options ...OperationOption,
) (*Result, error) {
	cancel := false

	loadedOptions := newOperationOptions(options...)

	var operationID uint32

	status := d.ffiMap.Cli.SendPromptedInput(
		d.ptr,
		&operationID,
		&cancel,
		input,
		prompt,
		loadedOptions.promptPattern,
		response,
		loadedOptions.hiddenInput,
		loadedOptions.abortInput,
		loadedOptions.requestedMode,
		string(loadedOptions.inputHandling),
		loadedOptions.retainTrailingPrompt,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit sendPromptedInput operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
