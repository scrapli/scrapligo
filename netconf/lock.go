package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newLockOptions(options ...Option) *lockOptions {
	o := &lockOptions{
		target: DatastoreTypeRunning,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type lockOptions struct {
	target DatastoreType
}

// Lock executes a netconf lock rpc. Supported options:
//   - WithTargetType
func (n *Netconf) Lock(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newLockOptions(options...)

	status := n.ffiMap.Netconf.Lock(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.target.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit lock operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
