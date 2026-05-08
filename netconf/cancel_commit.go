package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newCancelCommitOptions(options ...Option) *cancelCommitOptions {
	o := &cancelCommitOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type cancelCommitOptions struct {
	persistID string
}

// CancelCommit executes a netconf cancel-commit rpc.
func (n *Netconf) CancelCommit(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	_ = options

	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newCancelCommitOptions(options...)

	err := n.ffiMap.Netconf.CancelCommit(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.persistID,
	)
	if err != nil {
		return nil, err
	}

	return n.getResult(ctx, &cancel, operationID)
}
