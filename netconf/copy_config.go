package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newCopyConfigOptions(options ...Option) *copyConfigOptions {
	o := &copyConfigOptions{
		target: DatastoreTypeRunning,
		source: DatastoreTypeStartup,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type copyConfigOptions struct {
	target DatastoreType
	source DatastoreType
}

// CopyConfig executes a netconf copy config rpc. Supported options:
//   - WithTargetType
//   - WithSourceType
func (n *Netconf) CopyConfig(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newCopyConfigOptions(options...)

	err := n.ffiMap.Netconf.CopyConfig(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.target.String(),
		loadedOptions.source.String(),
	)
	if err != nil {
		return nil, err
	}

	return n.getResult(ctx, &cancel, operationID)
}
