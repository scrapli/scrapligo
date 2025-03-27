package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newEstablishSubscriptionOptions(options ...Option) *establishSubscriptionOptions {
	o := &establishSubscriptionOptions{
		stream:     DefaultStreamValue,
		filterType: FilterTypeSubtree,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type establishSubscriptionOptions struct {
	stream                string
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

// EstablishSubscription executes a netconf establish-subscription rpc. Supported options:
//   - WithStreamValue
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
func (d *Driver) EstablishSubscription(
	ctx context.Context,
	options ...Option,
) (*Result, error) {
	cancel := false

	var operationID uint32

	loadedOptions := newEstablishSubscriptionOptions(options...)

	// TODO ideally these are ordered how they used to be (and how the options are in docstring)
	//  just from a sanity perspective, buuuuut, that caused all sorts of issues w/ zig not "seeing"
	//  values in the right place... due to alignment maybe?
	status := d.ffiMap.Netconf.EstablishSubscription(
		d.ptr,
		&operationID,
		&cancel,
		loadedOptions.period,
		loadedOptions.stopTime,
		loadedOptions.dscp,
		loadedOptions.weighting,
		loadedOptions.dependency,
		loadedOptions.stream,
		loadedOptions.filter,
		loadedOptions.filterType.String(),
		loadedOptions.filterNamespacePrefix,
		loadedOptions.filterNamespace,
		loadedOptions.encoding,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError(
			"failed to submit establish-subscription operation",
			nil,
		)
	}

	return d.getResult(ctx, &cancel, operationID)
}
