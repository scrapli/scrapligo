package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// CancelCommit executes a netconf cancel-commit rpc.
func (d *Driver) CancelCommit(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	_ = options

	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := d.ffiMap.Netconf.CancelCommit(
		d.ptr,
		&operationID,
		&cancel,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit cancel-commit operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
