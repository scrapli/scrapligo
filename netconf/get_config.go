package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// GetConfigOption defines a functional option for a getconfig rpc.
type GetConfigOption func(o *getConfigOption)

func newGetConfigOptions(options ...GetConfigOption) getConfigOption {
	return getConfigOption{}
}

type getConfigOption struct{}

// GetConfig executes a netconf getconfig rpc.
func (n *Netconf) GetConfig(
	ctx context.Context,
	options ...GetConfigOption,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newGetConfigOptions(options...)
	_ = loadedOptions

	status := n.ffiMap.Netconf.GetConfig(
		n.ptr,
		&operationID,
		&cancel,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit getConfig operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
