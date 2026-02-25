package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newUnlockOptions(options ...Option) *unlockOptions {
	o := &unlockOptions{
		target: DatastoreTypeRunning,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type unlockOptions struct {
	target DatastoreType
}

// Unlock executes a netconf unlock rpc. Supported options:
//   - WithTargetType
func (n *Netconf) Unlock(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newUnlockOptions(options...)

	status := n.ffiMap.Netconf.Unlock(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.target.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit unlock operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
