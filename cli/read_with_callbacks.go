package cli

import (
	"bytes"
	"context"
	"fmt"
	"time"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligoutil "github.com/scrapli/scrapligo/util"
)

// ReadCallback represents a callback and how it is triggered, for use with ReadWithCallbacks.
type ReadCallback struct {
	name            string
	contains        string
	containsPattern string
	notContains     string
	once            bool
	completes       bool
	callback        func(c *Cli) error
}

// NewReadCallback returns a new ReadCallback with the given options set.
func NewReadCallback(
	name string,
	callback func(c *Cli) error,
	options ...Option,
) *ReadCallback {
	cb := &ReadCallback{
		name:     name,
		callback: callback,
	}

	for _, opt := range options {
		opt(cb)
	}

	return cb
}

func (r *ReadCallback) ok() bool {
	if r.contains == "" && r.containsPattern == "" {
		return false
	}

	return true
}

// ReadWithCallbacks optionally sends an initial "input" to the device and then continually reads
// from the session, checking new session output against the provided callbacks (in the order
// provided). When a callback is triggered, the callback is executed, if the callback is marked as
// "completes" then the parent function exits, otherwise this continues forever.
func (c *Cli) ReadWithCallbacks( //nolint: gocyclo
	ctx context.Context,
	initialInput string,
	callbacks ...*ReadCallback,
) (*Result, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	// check to make sure all callbacks have a contains or containspattern set
	for _, cb := range callbacks {
		if !cb.ok() {
			return nil, scrapligoerrors.NewOptionsError(
				fmt.Sprintf(
					"callback %q missing contains or containsPattern, cannot proceed", cb.name,
				), nil)
		}
	}

	startTime := time.Now().UnixNano()

	cancel := false

	if initialInput != "" {
		err := c.WriteAndReturn(initialInput)
		if err != nil {
			return nil, err
		}
	}

	pos := 0

	results := ""
	resultsRaw := bytes.NewBuffer(nil)

	executedCallbacks := make(map[string]any)

	for {
		var operationID uint32

		status := c.ffiMap.Cli.ReadAny(
			c.ptr,
			&operationID,
			&cancel,
		)
		if status != 0 {
			return nil, scrapligoerrors.NewFfiError("failed to submit readAny operation", nil)
		}

		r, err := c.getResult(ctx, &cancel, operationID)
		if err != nil {
			return nil, err
		}

		results += r.Result()
		resultsRaw.Write(r.ResultRaw())

		for _, cb := range callbacks {
			_, alreadyExecuted := executedCallbacks[cb.name]
			if alreadyExecuted && cb.once {
				continue
			}

			shouldExecute := false

			status = c.ffiMap.Cli.ReadCallbackShouldExecute(
				results[pos:],
				cb.name,
				cb.contains,
				cb.containsPattern,
				cb.notContains,
				&shouldExecute,
			)
			if status != 0 {
				return nil, scrapligoerrors.NewFfiError(
					"failed checking if callback should execute", nil,
				)
			}

			if !shouldExecute {
				continue
			}

			executedCallbacks[cb.name] = nil

			pos = len(results)

			err = cb.callback(c)
			if err != nil {
				return nil, err
			}

			if cb.completes {
				return NewResult(
					c.host,
					*c.options.Port,
					[]byte(initialInput),
					scrapligoutil.SafeInt64ToUint64(startTime),
					[]uint64{scrapligoutil.SafeInt64ToUint64(time.Now().UnixNano())},
					resultsRaw.Bytes(),
					[]byte(results),
					nil,
				), nil
			}
		}
	}
}
