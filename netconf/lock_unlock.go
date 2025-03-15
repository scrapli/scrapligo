package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newLockUnlockOptions(options ...Option) *lockUnlockOptions {
	o := &lockUnlockOptions{
		target: DatastoreTypeRunning,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type lockUnlockOptions struct {
	target DatastoreType
}

// Lock executes a netconf lock rpc. Supported options:
//   - WithTargetType
func (n *Netconf) Lock(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newLockUnlockOptions(options...)

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

// Unlock executes a netconf unlock rpc. Supported options:
//   - WithTargetType
func (n *Netconf) Unlock(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newLockUnlockOptions(options...)

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
