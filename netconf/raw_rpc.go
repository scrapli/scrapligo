package netconf

import (
	"context"
	"strings"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newRawRPCOptions(options ...Option) *rawRPCOptions {
	o := &rawRPCOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type rawRPCOptions struct {
	baseNamespacePrefix string
	extraNamespaces     [][2]string
}

func (o *rawRPCOptions) extraNamespacesToFFI() ([]byte, []uint64) { //nolint: gocritic
	namespaces := make([]string, len(o.extraNamespaces)*2) //nolint: mnd

	namespaceLens := make([]uint64, len(o.extraNamespaces)*2) //nolint: mnd

	var idx int

	for _, ns := range o.extraNamespaces {
		namespaces[idx] = ns[0]
		namespaceLens[idx] = uint64(len(ns[0]))

		idx++

		namespaces[idx] = ns[1]
		namespaceLens[idx] = uint64(len(ns[1]))

		idx++
	}

	return []byte(strings.Join(namespaces, "")), namespaceLens
}

// RawRPC executes a user provided "raw" rpc.
func (n *Netconf) RawRPC(
	ctx context.Context,
	payload string,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newRawRPCOptions(options...)

	extraNamespaces, extraNamespacesLens := loadedOptions.extraNamespacesToFFI()

	err := n.ffiMap.Netconf.RawRPC(
		n.ptr,
		&operationID,
		&cancel,
		payload,
		loadedOptions.baseNamespacePrefix,
		&extraNamespaces,
		&extraNamespacesLens,
	)
	if err != nil {
		return nil, err
	}

	return n.getResult(ctx, &cancel, operationID)
}
