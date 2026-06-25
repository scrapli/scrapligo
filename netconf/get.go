package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newGetOptions(options ...Option) *getOptions {
	o := &getOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type getOptions struct {
	filter                string
	filterType            *FilterType
	filterNamespacePrefix string
	filterNamespace       string
	defaultsType          *DefaultsType
}

func (o *getOptions) getFilterType() *uint8 {
	if o.filterType == nil {
		return nil
	}

	v := uint8(*o.filterType)

	return &v
}

func (o *getOptions) getDefaultsType() *uint8 {
	if o.defaultsType == nil {
		return nil
	}

	v := uint8(*o.defaultsType)

	return &v
}

// Get executes a netconf get rpc. Supported options:
//   - WithFilter
//   - WithFilterType
//   - WithFilterNamespacePrefix
//   - WithFilterNamespace
//   - WithDefaultsType
func (n *Netconf) Get(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newGetOptions(options...)

	err := n.ffiMap.Netconf.Get(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.filter,
		loadedOptions.getFilterType(),
		loadedOptions.filterNamespacePrefix,
		loadedOptions.filterNamespace,
		loadedOptions.getDefaultsType(),
	)
	if err != nil {
		return nil, err
	}

	return n.getResult(ctx, &cancel, operationID)
}
