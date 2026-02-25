package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

// Action executes a netconf action rpc.
func (n *Netconf) Action(
	ctx context.Context,
	action string,
	options ...Option,
) (*Result, error) {
	_ = options

	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := n.ffiMap.Netconf.Action(
		n.ptr,
		&operationID,
		&cancel,
		action,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit action operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
