package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// KillSession executes a netconf kill session rpc.
func (d *Driver) KillSession(
	ctx context.Context,
	sessionID uint64,
) (*Result, error) {
	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := d.ffiMap.Netconf.KillSession(
		d.ptr,
		&operationID,
		&cancel,
		sessionID,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit killSession operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
