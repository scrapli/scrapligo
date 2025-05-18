package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// Commit executes a netconf commit rpc.
func (n *Netconf) Commit(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	_ = options

	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := n.ffiMap.Netconf.Commit(
		n.ptr,
		&operationID,
		&cancel,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit commit operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
