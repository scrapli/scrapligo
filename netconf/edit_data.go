package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newEditDataOptions(options ...Option) *editDataOptions {
	o := &editDataOptions{
		datastore: DatastoreTypeRunning,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type editDataOptions struct {
	datastore DatastoreType
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

	status := n.ffiMap.Netconf.EditData(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.datastore.String(),
		content,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit edit-data operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
