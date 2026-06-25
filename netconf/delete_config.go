package netconf //nolint: dupl

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newDeleteConfigOptions(options ...Option) *deleteConfigOptions {
	o := &deleteConfigOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type deleteConfigOptions struct {
	target *DatastoreType
}

func (o *deleteConfigOptions) getTarget() *uint8 {
	if o.target == nil {
		return nil
	}

	v := uint8(*o.target)

	return &v
}

// DeleteConfig executes a netconf delete config rpc. Supported options:
//   - WithTargetType
func (n *Netconf) DeleteConfig(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newDeleteConfigOptions(options...)

	err := n.ffiMap.Netconf.DeleteConfig(
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
