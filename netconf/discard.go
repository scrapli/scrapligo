package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// Discard executes a netconf discard rpc.
func (d *Driver) Discard(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	_ = options

	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := d.ffiMap.Netconf.Discard(
		d.ptr,
		&operationID,
		&cancel,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit discard operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
