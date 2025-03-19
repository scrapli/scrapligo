package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
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
func (d *Driver) Unlock(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newUnlockOptions(options...)

	status := d.ffiMap.Netconf.Unlock(
		d.ptr,
		&operationID,
		&cancel,
		loadedOptions.target.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit unlock operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
