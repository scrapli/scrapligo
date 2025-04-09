package cli

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// SendInput sends an "input" to the device. Historically scrapli(go) had "SendCommand(s)" and
// "SendConfig(s)" operations, but these no longer exist. Instead, we have SendInput or SendInputs
// which accept their respective options -- the options can (among other things) control the "mode"
// (historically "privilege level") at which to send the input(s).
func (d *Driver) SendInput(
	ctx context.Context,
	input string,
	options ...OperationOption,
) (*Result, error) {
	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newOperationOptions(options...)

	status := d.ffiMap.Cli.SendInput(
		d.ptr,
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

	return d.getResult(ctx, &cancel, operationID)
}
