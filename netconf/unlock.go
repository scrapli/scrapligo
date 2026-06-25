package netconf //nolint: dupl

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newUnlockOptions(options ...Option) *unlockOptions {
	o := &unlockOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type unlockOptions struct {
	target *DatastoreType
}

func (o *unlockOptions) getTarget() *uint8 {
	if o.target == nil {
		return nil
	}

	v := uint8(*o.target)

	return &v
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

	err := n.ffiMap.Netconf.Unlock(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.getTarget(),
	)
	if err != nil {
		return nil, err
	}

	return n.getResult(ctx, &cancel, operationID)
}
