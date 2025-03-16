package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newGetOptions(options ...Option) *getOptions {
	o := &getOptions{
		filterType:   FilterTypeSubtree,
		defaultsType: DefaultsTypeUnset,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type getOptions struct {
	filter                string
	filterType            FilterType
	filterNamespacePrefix string
	filterNamespace       string
	defaultsType          DefaultsType
}

// Get executes a netconf get rpc. Supported options:
//   - WithFilter
//   - WithFilterType
//   - WithFilterNamespacePrefix
//   - WithFilterNamespace
//   - WithDefaultsType
func (d *Driver) Get(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newGetOptions(options...)

	status := d.ffiMap.Netconf.Get(
		d.ptr,
		&operationID,
		&cancel,
		loadedOptions.filter,
		loadedOptions.filterType.String(),
		loadedOptions.filterNamespacePrefix,
		loadedOptions.filterNamespace,
		loadedOptions.defaultsType.String(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit get operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
