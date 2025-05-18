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
func (n *Netconf) EditConfig(
	ctx context.Context,
	config string,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newEditConfigOptions(options...)

	status := n.ffiMap.Netconf.EditConfig(
		n.ptr,
		&operationID,
		&cancel,
		config,
		loadedOptions.target.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit editConfig operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
