package cli

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// GetPrompt returns a Result object containing the current "prompt" of the target device.
func (d *Driver) GetPrompt(ctx context.Context) (*Result, error) {
	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := d.ffiMap.Cli.GetPrompt(d.ptr, &operationID, &cancel)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit getPrompt operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
