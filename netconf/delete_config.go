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
func (n *Netconf) DeleteConfig(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newDeleteConfigOptions(options...)

	status := n.ffiMap.Netconf.DeleteConfig(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.target.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit delete-config operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
