package network

import (
	"errors"

	"github.com/scrapli/scrapligo/util"
)

const (
	defaultPrivilegeLevel = "configuration"
)

// OperationOptions is a struct containing "operation" options that are relevant to the network
// Driver, for example providing a target privilege level for a SendInteractive operation.
type OperationOptions struct {
	PrivilegeLevel string
}

// NewOperation returns a new OperationOptions object with the defaults set and any provided options
// applied.
func NewOperation(options ...util.Option) (*OperationOptions, error) {
	o := &OperationOptions{
		PrivilegeLevel: defaultPrivilegeLevel,
	}

	for _, option := range options {
		err := option(o)
		if err != nil {
			if !errors.Is(err, util.ErrIgnoredOption) {
				return nil, err
			}
		}
	}

	if o.PrivilegeLevel == "" {
		o.PrivilegeLevel = defaultPrivilegeLevel
	}

	return o, nil
}
