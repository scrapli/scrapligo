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

	status := n.ffiMap.Netconf.Validate(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.source.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit validate operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
