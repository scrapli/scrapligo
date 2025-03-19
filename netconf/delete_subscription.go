package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// DeleteSubscription executes a netconf delete-subscription rpc.
func (d *Driver) DeleteSubscription(
	ctx context.Context,
	id uint64,
	options ...Option,
) (*Result, error) {
	_ = options

	cancel := false

	var operationID uint32

	status := d.ffiMap.Netconf.DeleteSubscription(
		d.ptr,
		&operationID,
		&cancel,
		id,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError(
			"failed to submit delete-subscription operation",
			nil,
		)
	}

	return d.getResult(ctx, &cancel, operationID)
}
