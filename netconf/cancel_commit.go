package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// CancelCommit executes a netconf cancel-commit rpc.
func (n *Netconf) CancelCommit(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	_ = options

	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := n.ffiMap.Netconf.CancelCommit(
		n.ptr,
		&operationID,
		&cancel,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit cancel-commit operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
