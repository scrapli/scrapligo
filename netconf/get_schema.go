package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newGetSchemaOptions(options ...Option) *getSchemaOptions {
	o := &getSchemaOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type getSchemaOptions struct {
	version string
	format  *SchemaFormat
}

func (o *getSchemaOptions) getFormat() *uint8 {
	if o.format == nil {
		return nil
	}

	v := uint8(*o.format)

	return &v
}

// GetSchema executes a netconf get-schema rpc
//   - WithSchemaFormat
//   - WithVersion
func (n *Netconf) GetSchema(
	ctx context.Context,
	identifier string,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newGetSchemaOptions(options...)

	err := n.ffiMap.Netconf.GetSchema(
		n.ptr,
		&operationID,
		&cancel,
		identifier,
		loadedOptions.version,
		loadedOptions.getFormat(),
	)
	if err != nil {
		return nil, err
	}

	return n.getResult(ctx, &cancel, operationID)
}
