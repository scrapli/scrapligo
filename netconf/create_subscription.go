package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newCreateSubscriptionOptions(options ...Option) *createSubscriptionOptions {
	o := &createSubscriptionOptions{
		filterType: FilterTypeSubtree,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type createSubscriptionOptions struct {
	filter                string
	filterType            FilterType
	filterNamespacePrefix string
	filterNamespace       string
	startTime             uint64
	stopTime              uint64
}

// CreateSubscription executes a netconf create-subscription rpc. Supported options:
//   - WithFilter
//   - WithFilterType
//   - WithFilterNamespacePrefix
//   - WithFilterNamespace
//   - WithStartTime
//   - WithStopTime
func (d *Driver) CreateSubscription(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newCreateSubscriptionOptions(options...)

	status := d.ffiMap.Netconf.CreateSubscription(
		d.ptr,
		&operationID,
		&cancel,
		"stream", // TODO -- also are the start/stop optional or is at least one of those required?
		loadedOptions.filter,
		loadedOptions.filterType.String(),
		loadedOptions.filterNamespacePrefix,
		loadedOptions.filterNamespace,
		loadedOptions.startTime,
		loadedOptions.stopTime,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError(
			"failed to submit create-subscription operation",
			nil,
		)
	}

	return d.getResult(ctx, &cancel, operationID)
}
