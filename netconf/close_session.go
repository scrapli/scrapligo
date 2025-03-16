package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// CloseSession executes a netconf delete config rpc. Supported options:
//   - WithTargetType
func (d *Driver) CloseSession(
	ctx context.Context,
) (*Result, error) {
	cancel := false

	var operationID uint32

	status := d.ffiMap.Netconf.CloseSession(
		d.ptr,
		&operationID,
		&cancel,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit closeSession operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
