package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
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
func (n *Netconf) GetData(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newGetDataOptions(options...)

	status := n.ffiMap.Netconf.GetData(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.datastore.String(),
		loadedOptions.filter,
		loadedOptions.filterType.String(),
		loadedOptions.filterNamespacePrefix,
		loadedOptions.filterNamespace,
		loadedOptions.configFilter.String(),
		loadedOptions.originFilters,
		loadedOptions.maxDepth,
		loadedOptions.withOrigin,
		loadedOptions.defaultsType.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit get-data operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
