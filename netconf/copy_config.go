package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newCopyConfigOptions(options ...Option) *copyConfigOptions {
	o := &copyConfigOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type copyConfigOptions struct {
	target *DatastoreType
	source *DatastoreType
}

func (o *copyConfigOptions) getTarget() *uint8 {
	if o.target == nil {
		return nil
	}

	v := uint8(*o.target)

	return &v
}

func (o *copyConfigOptions) getSource() *uint8 {
	if o.source == nil {
		return nil
	}

	v := uint8(*o.source)

	return &v
}

// CopyConfig executes a netconf copy config rpc. Supported options:
//   - WithTargetType
//   - WithSourceType
func (n *Netconf) CopyConfig(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newCopyConfigOptions(options...)

	err := n.ffiMap.Netconf.CopyConfig(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.getTarget(),
		loadedOptions.getSource(),
	)
	if err != nil {
		return nil, err
	}

	return n.getResult(ctx, &cancel, operationID)
}
