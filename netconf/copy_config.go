package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
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
	cancel := false

	var operationID uint32

	loadedOptions := newCopyConfigOptions(options...)

	status := n.ffiMap.Netconf.CopyConfig(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.target.String(),
		loadedOptions.source.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit copyConfig operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
