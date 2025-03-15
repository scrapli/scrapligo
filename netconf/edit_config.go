package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newEditConfigOptions(options ...Option) *editConfigOption {
	o := &editConfigOption{
		target: DatastoreTypeRunning,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type editConfigOption struct {
	target DatastoreType
}

// EditConfig executes a netconf getconfig rpc. Supported options:
//   - WithTargetType
func (n *Netconf) EditConfig(
	ctx context.Context,
	config string,
	options ...Option,
) (*Result, error) {
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
