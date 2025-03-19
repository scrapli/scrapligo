package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newModifySubscriptionOptions(options ...Option) *modifySubscriptionOptions {
	o := &modifySubscriptionOptions{
		filterType: FilterTypeSubtree,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type modifySubscriptionOptions struct {
	filter                string
	filterType            FilterType
	filterNamespacePrefix string
	filterNamespace       string
	period                uint64
	stopTime              uint64
	dscp                  uint8
	weighting             uint8
	dependency            uint32
	encoding              string
}

// ModifySubscription executes a netconf modify-subscription rpc. Supported options:
//   - WithFilter
//   - WithFilterType
//   - WithFilterNamespacePrefix
//   - WithFilterNamespace
//   - WithPeriod
//   - WithStopTime
//   - WithDSCP
//   - WithWeighting
//   - WithDependency
//   - WithEncoding
func (d *Driver) ModifySubscription(
	ctx context.Context,
	id uint64,
	options ...Option,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newModifySubscriptionOptions(options...)

	status := d.ffiMap.Netconf.ModifySubscription(
		d.ptr,
		&operationID,
		&cancel,
		id,
		"stream", // TODO -- also are the start/stop optional or is at least one of those required?
		loadedOptions.filter,
		loadedOptions.filterType.String(),
		loadedOptions.filterNamespacePrefix,
		loadedOptions.filterNamespace,
		loadedOptions.period,
		loadedOptions.stopTime,
		loadedOptions.dscp,
		loadedOptions.weighting,
		loadedOptions.dependency,
		loadedOptions.encoding,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError(
			"failed to submit modify-subscription operation",
			nil,
		)
	}

	return d.getResult(ctx, &cancel, operationID)
}
