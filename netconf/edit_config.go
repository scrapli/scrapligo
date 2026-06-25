package netconf

import (
	"context"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

func newEditConfigOptions(options ...Option) *editConfigOptions {
	o := &editConfigOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type editConfigOptions struct {
	target           *DatastoreType
	defaultOperation *DefaultOperation
	testOption       *TestOption
	errorOption      *ErrorOption
}

func (o *editConfigOptions) getTarget() *uint8 {
	if o.target == nil {
		return nil
	}

	v := uint8(*o.target)

	return &v
}

func (o *editConfigOptions) getDefaultOperation() *uint8 {
	if o.defaultOperation == nil {
		return nil
	}

	v := uint8(*o.defaultOperation)

	return &v
}

func (o *editConfigOptions) getTestOption() *uint8 {
	if o.testOption == nil {
		return nil
	}

	v := uint8(*o.testOption)

	return &v
}

func (o *editConfigOptions) getErrorOption() *uint8 {
	if o.errorOption == nil {
		return nil
	}

	v := uint8(*o.errorOption)

	return &v
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

	err := n.ffiMap.Netconf.EditConfig(
		n.ptr,
		&operationID,
		&cancel,
		config,
		loadedOptions.getTarget(),
		loadedOptions.getDefaultOperation(),
		loadedOptions.getTestOption(),
		loadedOptions.getErrorOption(),
	)
	if err != nil {
		return nil, err
	}

	return n.getResult(ctx, &cancel, operationID)
}
