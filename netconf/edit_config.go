package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newEditConfigOptions(options ...Option) *editConfigOptions {
	o := &editConfigOptions{
		target: DatastoreTypeRunning,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type editConfigOptions struct {
	target           DatastoreType
	defaultOperation *DefaultOperation
	testOption       *TestOption
	errorOption      *ErrorOption
}

func (o *editConfigOptions) getDefaultOperation() string {
	if o.defaultOperation == nil {
		return ""
	}

	return o.defaultOperation.String()
}

func (o *editConfigOptions) getTestOption() string {
	if o.testOption == nil {
		return ""
	}

	return o.testOption.String()
}

func (o *editConfigOptions) getErrorOption() string {
	if o.errorOption == nil {
		return ""
	}

	return o.errorOption.String()
}

// EditConfig executes a netconf edit config rpc. Supported options:
//   - WithTargetType
func (n *Netconf) EditConfig(
	ctx context.Context,
	config string,
	options ...Option,
) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newEditConfigOptions(options...)

	status := n.ffiMap.Netconf.EditConfig(
		n.ptr,
		&operationID,
		&cancel,
		config,
		loadedOptions.target.String(),
		loadedOptions.getDefaultOperation(),
		loadedOptions.getTestOption(),
		loadedOptions.getErrorOption(),
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit editConfig operation", nil)
	}

	return n.getResult(ctx, &cancel, operationID)
}
