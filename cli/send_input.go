package cli

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newSendInputOptions(options ...Option) *sendInputOptions {
	o := &sendInputOptions{
		inputHandling: InputHandlingFuzzy,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type sendInputOptions struct {
	requestedMode        string
	inputHandling        InputHandling
	retainInput          bool
	retainTrailingPrompt bool
}

// SendInput sends an "input" to the device. Historically scrapli(go) had "SendCommand(s)" and
// "SendConfig(s)" operations, but these no longer exist. Instead, we have SendInput or SendInputs
// which accept their respective options -- the options can (among other things) control the "mode"
// (historically "privilege level") at which to send the input(s).
func (c *Cli) SendInput(
	ctx context.Context,
	input string,
	options ...Option,
) (*Result, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newSendInputOptions(options...)

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

	return c.getResult(ctx, &cancel, operationID)
}
