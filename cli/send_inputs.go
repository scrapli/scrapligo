package cli

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// SendInputs send multiple "inputs" to the device.
func (d *Driver) SendInputs(
	ctx context.Context,
	inputs []string,
	options ...OperationOption,
) (*MultiResult, error) {
	cancel := false

	loadedOptions := newOperationOptions(options...)

	results := NewMultiResult(d.host, *d.options.Port)

	for _, input := range inputs {
		var operationID uint32

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

		r, err := d.getResult(ctx, &cancel, operationID)
		if err != nil {
			return nil, err
		}

		results.AppendResult(r)

		if r.Failed && loadedOptions.stopOnIndicatedFailure {
			// note that this returns nil for an error since there was nothing unrecoverable
			// (probably) that happened, just we saw some stuff in the output saying that we
			// had a bad input or something
			return results, nil
		}
	}

	return results, nil
}
