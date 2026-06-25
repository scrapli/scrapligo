package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newGetDataOptions(options ...Option) *getDataOptions {
	o := &getDataOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type getDataOptions struct {
	datastore             *DatastoreType
	filter                string
	filterType            *FilterType
	filterNamespacePrefix string
	filterNamespace       string
	configFilter          *ConfigFilter
	originFilters         string
	maxDepth              uint32
	withOrigin            bool
	defaultsType          *DefaultsType
}

func (o *getDataOptions) getDatastore() *uint8 {
	if o.datastore == nil {
		return nil
	}

	v := uint8(*o.datastore)

	return &v
}

func (o *getDataOptions) getFilterType() *uint8 {
	if o.filterType == nil {
		return nil
	}

	v := uint8(*o.filterType)

	return &v
}

func (o *getDataOptions) getConfigFilter() *bool {
	if o.configFilter == nil {
		return nil
	}

	var v bool

	switch *o.configFilter {
	case ConfigFilterTrue:
		v = true
	case ConfigFilterFalse:
		v = false
	}

	return &v
}

func (o *getDataOptions) getDefaultsType() *uint8 {
	if o.defaultsType == nil {
		return nil
	}

	v := uint8(*o.defaultsType)

	return &v
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

	err := n.ffiMap.Netconf.GetData(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.getDatastore(),
		loadedOptions.filter,
		loadedOptions.getFilterType(),
		loadedOptions.filterNamespacePrefix,
		loadedOptions.filterNamespace,
		loadedOptions.getConfigFilter(),
		loadedOptions.originFilters,
		loadedOptions.maxDepth,
		loadedOptions.withOrigin,
		loadedOptions.getDefaultsType(),
	)
	if err != nil {
		return nil, err
	}

	return n.getResult(ctx, &cancel, operationID)
}
