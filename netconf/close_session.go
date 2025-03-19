package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// CloseSession executes a netconf close-session rpc.
func (d *Driver) CloseSession(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	_ = options

	cancel := false

	var operationID uint32

	status := d.ffiMap.Netconf.CloseSession(
		d.ptr,
		&operationID,
		&cancel,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit close-session operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
