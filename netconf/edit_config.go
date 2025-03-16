package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newEditConfigOptions(options ...Option) *editConfigOptions {
	o := &editConfigOptions{
		target: DatastoreTypeRunning,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type editConfigOptions struct {
	target DatastoreType
}

// EditConfig executes a netconf edit config rpc. Supported options:
//   - WithTargetType
func (d *Driver) EditConfig(
	ctx context.Context,
	config string,
	options ...Option,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newEditConfigOptions(options...)

	status := d.ffiMap.Netconf.EditConfig(
		d.ptr,
		&operationID,
		&cancel,
		config,
		loadedOptions.target.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit editConfig operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
