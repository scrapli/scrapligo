package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newGetConfigOptions(options ...Option) *getConfigOptions {
	o := &getConfigOptions{
		source:       DatastoreTypeRunning,
		filterType:   FilterTypeSubtree,
		defaultsType: DefaultsTypeUnset,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type getConfigOptions struct {
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
func (d *Driver) GetConfig(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newGetConfigOptions(options...)

	status := d.ffiMap.Netconf.GetConfig(
		d.ptr,
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

	return d.getResult(ctx, &cancel, operationID)
}
