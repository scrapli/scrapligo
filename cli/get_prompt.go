package cli

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

// GetPrompt returns a Result object containing the current "prompt" of the target device.
func (c *Cli) GetPrompt(ctx context.Context) (*Result, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	err := c.ffiMap.Cli.GetPrompt(c.ptr, &operationID, &cancel)
	if err != nil {
		return nil, err
	}

	return c.getResult(ctx, &cancel, operationID)
}
