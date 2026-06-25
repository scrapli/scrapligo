package netconf //nolint: dupl

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newValidateOptions(options ...Option) *validateOptions {
	o := &validateOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type validateOptions struct {
	source *DatastoreType
}

func (o *validateOptions) getSource() *uint8 {
	if o.source == nil {
		return nil
	}

	v := uint8(*o.source)

	return &v
}

// Validate executes a netconf validate rpc. Supported options:
//   - WithSource
func (n *Netconf) Validate(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newValidateOptions(options...)

	err := n.ffiMap.Netconf.Validate(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.getSource(),
	)
	if err != nil {
		return nil, err
	}

	return n.getResult(ctx, &cancel, operationID)
}
