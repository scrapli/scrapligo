package cli

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// EnterMode is used to explicitly enter a mode (i.e. enter "config mode" or "shell" or some other
// platform specific "mode").
func (d *Driver) EnterMode(ctx context.Context, requestedMode string) (*Result, error) {
	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := d.ffiMap.Cli.EnterMode(d.ptr, &operationID, &cancel, requestedMode)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit enterMode operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
