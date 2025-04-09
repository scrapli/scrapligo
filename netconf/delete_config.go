package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newDeleteConfigOptions(options ...Option) *deleteConfigOptions {
	o := &deleteConfigOptions{
		target: DatastoreTypeRunning,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type deleteConfigOptions struct {
	target DatastoreType
}

// DeleteConfig executes a netconf delete config rpc. Supported options:
//   - WithTargetType
func (d *Driver) DeleteConfig(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newDeleteConfigOptions(options...)

	status := d.ffiMap.Netconf.DeleteConfig(
		d.ptr,
		&operationID,
		&cancel,
		loadedOptions.target.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit delete-config operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
