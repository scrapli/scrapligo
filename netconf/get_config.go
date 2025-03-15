package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newGetConfigOptions(options ...Option) *getConfigOption {
	o := &getConfigOption{
		source:       DatastoreTypeRunning,
		filterType:   FilterTypeSubtree,
		defaultsType: DefaultsTypeUnset,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type getConfigOption struct {
	source                DatastoreType
	filter                string
	filterType            FilterType
	filterNamespacePrefix string
	filterNamespace       string
	defaultsType          DefaultsType
}

// GetConfig executes a netconf getconfig rpc. Supported options:
//   - WithSourceType
//   - WithFilter
//   - WithFilterType
//   - WithFilterNamespacePrefix
//   - WithFilterNamespace
//   - WithDefaultsType
func (n *Netconf) GetConfig(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newGetConfigOptions(options...)

	status := n.ffiMap.Netconf.GetConfig(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.source.String(),
		loadedOptions.filter,
		loadedOptions.filterType.String(),
		loadedOptions.filterNamespacePrefix,
		loadedOptions.filterNamespace,
		loadedOptions.defaultsType.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit getConfig operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
