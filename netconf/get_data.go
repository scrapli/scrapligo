package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newGetDataOptions(options ...Option) *getDataOptions {
	o := &getDataOptions{
		filterType:   FilterTypeSubtree,
		defaultsType: DefaultsTypeUnset,
		configFilter: ConfigFilterUnset,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type getDataOptions struct {
	datastore             DatastoreType
	filter                string
	filterType            FilterType
	filterNamespacePrefix string
	filterNamespace       string
	configFilter          ConfigFilter
	originFilters         string
	maxDepth              uint32
	withOrigin            bool
	defaultsType          DefaultsType
}

// GetData executes a netconf get-data rpc. Supported options:
//   - WithDatastore
//   - WithFilter
//   - WithFilterType
//   - WithFilterNamespacePrefix
//   - WithFilterNamespace
//   - WithDefaultsType
//   - WithConfigFilter
//   - WithMaxDepth
//   - WithOrigin
//   - WithDefaultsType
func (d *Driver) GetData(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newGetDataOptions(options...)

	status := d.ffiMap.Netconf.GetData(
		d.ptr,
		&operationID,
		&cancel,
		loadedOptions.configFilter.String(),
		loadedOptions.maxDepth,
		loadedOptions.withOrigin,
		loadedOptions.datastore.String(),
		loadedOptions.filter,
		loadedOptions.filterType.String(),
		loadedOptions.filterNamespacePrefix,
		loadedOptions.filterNamespace,
		loadedOptions.originFilters,
		loadedOptions.defaultsType.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit get-data operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
