package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// Commit executes a netconf commit rpc.
func (d *Driver) Commit(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	_ = options

	cancel := false

	var operationID uint32

	status := d.ffiMap.Netconf.Commit(
		d.ptr,
		&operationID,
		&cancel,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit commit operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
