package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// Action executes a netconf action rpc.
func (d *Driver) Action(
	ctx context.Context,
	action string,
	options ...Option,
) (*Result, error) {
	_ = options

	cancel := false

	var operationID uint32

	status := d.ffiMap.Netconf.Action(
		d.ptr,
		&operationID,
		&cancel,
		action,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit action operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
