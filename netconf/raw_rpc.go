package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// RawRPC executes a user provided "raw" rpc.
func (n *Netconf) RawRPC(
	ctx context.Context,
	payload string,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := n.ffiMap.Netconf.RawRPC(
		n.ptr,
		&operationID,
		&cancel,
		payload,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit raw-rpc operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
