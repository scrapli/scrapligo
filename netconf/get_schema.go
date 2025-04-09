package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newGetSchemaOptions(options ...Option) *getSchemaOptions {
	o := &getSchemaOptions{
		format: SchemaFormatYang,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type getSchemaOptions struct {
	version string
	format  SchemaFormat
}

// GetSchema executes a netconf get-schema rpc
//   - WithSchemaFormat
//   - WithVersion
func (d *Driver) GetSchema(
	ctx context.Context,
	identifier string,
	options ...Option,
) (*Result, error) {
	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newGetSchemaOptions(options...)

	status := d.ffiMap.Netconf.GetSchema(
		d.ptr,
		&operationID,
		&cancel,
		identifier,
		loadedOptions.version,
		loadedOptions.format.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit get-schema operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
