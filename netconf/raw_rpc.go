package netconf

import (
	"context"
	"fmt"
	"strings"

	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

func newRawRPCOptions(options ...Option) *rawRPCOptions {
	o := &rawRPCOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type rawRPCOptions struct {
	extraNamespaces [][2]string
}

func (o *rawRPCOptions) extraNamespacesToFFI() string {
	namespaces := make([]string, len(o.extraNamespaces))

	for i, ns := range o.extraNamespaces {
		namespaces[i] = fmt.Sprintf("%s::%s", ns[0], ns[1])
	}

	return strings.Join(namespaces, scrapligoconstants.LibScrapliDelimiter)
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

	status := n.ffiMap.Netconf.RawRPC(
		n.ptr,
		&operationID,
		&cancel,
		payload,
		loadedOptions.extraNamespacesToFFI(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit raw-rpc operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
