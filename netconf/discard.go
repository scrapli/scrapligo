package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// Discard executes a netconf discard rpc.
func (n *Netconf) Discard(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	_ = options

	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := n.ffiMap.Netconf.Discard(
		n.ptr,
		&operationID,
		&cancel,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit discard operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
