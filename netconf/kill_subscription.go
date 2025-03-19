package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// KillSubscription executes a netconf kill-subscription rpc.
func (d *Driver) KillSubscription(
	ctx context.Context,
	id uint64,
	options ...Option,
) (*Result, error) {
	_ = options

	cancel := false

	var operationID uint32

	status := d.ffiMap.Netconf.KillSubscription(
		d.ptr,
		&operationID,
		&cancel,
		id,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError(
			"failed to submit kill-subscription operation",
			nil,
		)
	}

	return d.getResult(ctx, &cancel, operationID)
}
