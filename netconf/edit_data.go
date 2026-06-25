package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newEditDataOptions(options ...Option) *editDataOptions {
	o := &editDataOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type editDataOptions struct {
	datastore        *DatastoreType
	defaultOperation *DefaultOperation
}

func (o *editDataOptions) getDatastore() *uint8 {
	if o.datastore == nil {
		return nil
	}

	v := uint8(*o.datastore)

	return &v
}

func (o *editDataOptions) getDefaultOperation() *uint8 {
	if o.defaultOperation == nil {
		return nil
	}

	v := uint8(*o.defaultOperation)

	return &v
}

// EditData executes a netconf edit-data rpc. Supported options:
//   - WithDatastore
func (n *Netconf) EditData(
	ctx context.Context,
	content string,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newEditDataOptions(options...)

	err := n.ffiMap.Netconf.EditData(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.getDatastore(),
		content,
		loadedOptions.getDefaultOperation(),
	)
	if err != nil {
		return nil, err
	}

	return n.getResult(ctx, &cancel, operationID)
}
