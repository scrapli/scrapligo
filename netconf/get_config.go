package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newGetConfigOptions(options ...Option) *getConfigOptions {
	o := &getConfigOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type getConfigOptions struct {
	source                *DatastoreType
	filter                string
	filterType            *FilterType
	filterNamespacePrefix string
	filterNamespace       string
	defaultsType          *DefaultsType
}

func (o *getConfigOptions) getSource() *uint8 {
	if o.source == nil {
		return nil
	}

	v := uint8(*o.source)

	return &v
}

func (o *getConfigOptions) getFilterType() *uint8 { return nil }

func (o *getConfigOptions) getDefaultsType() *uint8 { return nil }

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
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newGetConfigOptions(options...)

	err := n.ffiMap.Netconf.GetConfig(
		n.ptr,
		&operationID,
		&cancel,
		loadedOptions.getSource(),
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
