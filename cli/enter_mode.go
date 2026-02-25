package cli

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

// EnterMode is used to explicitly enter a mode (i.e. enter "config mode" or "shell" or some other
// platform specific "mode").
func (c *Cli) EnterMode(ctx context.Context, requestedMode string) (*Result, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := c.ffiMap.Cli.EnterMode(c.ptr, &operationID, &cancel, requestedMode)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit enterMode operation", nil)
	}

	return c.getResult(ctx, &cancel, operationID)
}
