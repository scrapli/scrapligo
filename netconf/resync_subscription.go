package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// ResyncSubscription executes a netconf resync-subscription rpc.
func (d *Driver) ResyncSubscription(
	ctx context.Context,
	id uint64,
	options ...Option,
) (*Result, error) {
	_ = options

	cancel := false

	var operationID uint32

	status := d.ffiMap.Netconf.ResyncSubscription(
		d.ptr,
		&operationID,
		&cancel,
		id,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError(
			"failed to submit resync-subscription operation",
			nil,
		)
	}

	return d.getResult(ctx, &cancel, operationID)
}
