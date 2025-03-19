package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newValidateOptions(options ...Option) *validateOptions {
	o := &validateOptions{
		source: DatastoreTypeRunning,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type validateOptions struct {
	source DatastoreType
}

// Validate executes a netconf validate rpc. Supported options:
//   - WithTargetType
func (d *Driver) Validate(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newValidateOptions(options...)

	status := d.ffiMap.Netconf.Validate(
		d.ptr,
		&operationID,
		&cancel,
		loadedOptions.source.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit validate operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
