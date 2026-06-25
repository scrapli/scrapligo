package netconf //nolint: dupl

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newLockOptions(options ...Option) *lockOptions {
	o := &lockOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type lockOptions struct {
	target *DatastoreType
}

func (o *lockOptions) getTarget() *uint8 {
	if o.target == nil {
		return nil
	}

	v := uint8(*o.target)

	return &v
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

	err := n.ffiMap.Netconf.Lock(
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
