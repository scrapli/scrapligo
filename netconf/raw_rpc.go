package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// RawRPC executes a user provided "raw" rpc.
func (d *Driver) RawRPC(
	ctx context.Context,
	payload string,
) (*Result, error) {
	cancel := false

	var operationID uint32

	status := d.ffiMap.Netconf.RawRPC(
		d.ptr,
		&operationID,
		&cancel,
		payload,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit raw-rpc operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
