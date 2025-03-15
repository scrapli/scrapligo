package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// KillSession executes a netconf kill session rpc.
func (n *Netconf) KillSession(
	ctx context.Context,
	sessionID uint64,
) (*Result, error) {
	cancel := false

	var operationID uint32

	status := n.ffiMap.Netconf.KillSession(
		n.ptr,
		&operationID,
		&cancel,
		sessionID,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit killSession operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
