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
func (d *Driver) EditData(
	ctx context.Context,
	content string,
	options ...Option,
) (*Result, error) {
	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newEditDataOptions(options...)

	status := d.ffiMap.Netconf.EditData(
		d.ptr,
		&operationID,
		&cancel,
		loadedOptions.datastore.String(),
		content,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit edit-data operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
