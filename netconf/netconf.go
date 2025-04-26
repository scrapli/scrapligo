package netconf

import (
	"context"
	"math/rand"
	"time"

	"github.com/ebitengine/purego"
	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
	scrapligointernal "github.com/scrapli/scrapligo/internal"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligoutil "github.com/scrapli/scrapligo/util"
)

// Netconf is an object representing a netconf connection to a device of some sort -- this object
// wraps the underlying zig (netconf) driver (created via libscrapli).
type Netconf struct {
	ptr     uintptr
	ffiMap  *scrapligoffi.Mapping
	host    string
	options *scrapligointernal.Options

	minPollDelay  uint64
	maxPollDelay  uint64
	backoffFactor uint8
}

// NewNetconf returns a new instance of Netconf setup with the given options.
func NewNetconf(
	host string,
	opts ...scrapligooptions.Option,
) (*Netconf, error) {
	ffiMap, err := scrapligoffi.GetMapping()
	if err != nil {
		return nil, err
	}

	n := &Netconf{
		ffiMap:  ffiMap,
		host:    host,
		options: scrapligointernal.NewOptions(),
	}

	for _, opt := range opts {
		err = opt(n.options)
		if err != nil {
			return nil, scrapligoerrors.NewOptionsError("failed applying option", err)
		}
	}

	if n.options.Port == nil {
		p := scrapligoconstants.DefaultNetconfPort

		n.options.Port = &p
	}

	n.minPollDelay = scrapligoconstants.DefaultReadDelayMinNs
	if n.options.Session.ReadDelayMinNs != nil {
		n.minPollDelay = *n.options.Session.ReadDelayMinNs
	}

	n.maxPollDelay = scrapligoconstants.DefaultReadDelayMaxNs
	if n.options.Session.ReadDelayMaxNs != nil {
		n.maxPollDelay = *n.options.Session.ReadDelayMaxNs
	}

	n.backoffFactor = scrapligoconstants.DefaultReadDelayBackoffFactor
	if n.options.Session.ReadDelayBackoffFactor != nil {
		n.backoffFactor = *n.options.Session.ReadDelayBackoffFactor
	}

	return n, nil
}

// GetPtr returns the pointer to the zig driver, don't use this unless you know what you are doing,
// this is just exposed so you *can* get to it if you want to.
func (n *Netconf) GetPtr() (uintptr, *scrapligoffi.Mapping) {
	return n.ptr, n.ffiMap
}

// Open opens the driver object. This method spawns the underlying zig driver which the Cli then
// holds a pointer to. All Cli operations operate against this pointer (though this is
// transparent to the user).
func (n *Netconf) Open(ctx context.Context) (*Result, error) {
	// ensure we dealloc if something happens, otherwise users calls to defer close would not be
	// super handy
	cleanup := false

	defer func() {
		if !cleanup {
			return
		}

		n.ffiMap.Shared.Free(n.ptr)
	}()

	var loggerCallback uintptr
	if n.options.LoggerCallback != nil {
		loggerCallback = purego.NewCallback(n.options.LoggerCallback)
	}

	n.ptr = n.ffiMap.Netconf.Alloc(
		loggerCallback,
		n.host,
		*n.options.Port,
		string(n.options.TransportKind),
	)

	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("failed to allocate driver", nil)
	}

	err := n.options.Apply(n.ptr, n.ffiMap)
	if err != nil {
		return nil, scrapligoerrors.NewFfiError("failed to applying driver options", err)
	}

	cancel := false

	var operationID uint32

	status := n.ffiMap.Shared.Open(n.ptr, &operationID, &cancel)
	if status != 0 {
		cleanup = true

		return nil, scrapligoerrors.NewFfiError("failed to submit open operation", nil)
	}

	result, err := n.getResult(ctx, &cancel, operationID)
	if err != nil {
		cleanup = true

		return nil, err
	}

	return result, nil
}

// Close closes the netconf object. This also deallocates the underlying (zig) netconf object.
func (n *Netconf) Close(ctx context.Context) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	status := n.ffiMap.Shared.Close(n.ptr, &operationID, &cancel)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit close operation", nil)
	}

	result, err := n.getResult(ctx, &cancel, operationID)

	n.ffiMap.Shared.Free(n.ptr)

	return result, err
}

// GetSessionID returns the session-id as parsed during the capabilities exchange -- if we for some
// reason didn't parse the session-id during capabilities exchange this will return an error.
func (n *Netconf) GetSessionID() (uint64, error) {
	if n.ptr == 0 {
		return 0, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	var sessionID uint64

	status := n.ffiMap.Netconf.GetSessionID(n.ptr, &sessionID)
	if status != 0 {
		return 0, scrapligoerrors.NewFfiError("session-id not set", nil)
	}

	return sessionID, nil
}

func getPollDelay(curVal, minVal, maxVal uint64, backoffFactor uint8) uint64 {
	newVal := curVal
	newVal *= uint64(backoffFactor)

	if newVal > maxVal {
		newVal = maxVal
	}

	if minVal == 0 {
		return newVal
	}

	return newVal + scrapligoutil.SafeInt64ToUint64(
		rand.Int63n(scrapligoutil.SafeUint64ToInt64(minVal)), //nolint:gosec
	)
}

func (n *Netconf) getResult(
	ctx context.Context,
	cancel *bool,
	operationID uint32,
) (*Result, error) {
	var done bool

	var inputSize, resultRawSize, resultSize, rpcWarningsSize, rpcErrorsSize, errSize uint64

	curPollDelay := scrapligoconstants.DefaultReadDelayMaxNs

	for {
		select {
		case <-ctx.Done():
			*cancel = true

			return nil, ctx.Err()
		default:
		}

		// we obviously cant have too tight a loop here or cpu will go nuts and we'll block things,
		// so we'll sleep the same as the zig read delay will be
		time.Sleep(time.Duration(scrapligoutil.SafeUint64ToInt64(curPollDelay)))

		rc := n.ffiMap.Netconf.PollOperation(
			n.ptr,
			operationID,
			&done,
			&inputSize,
			&resultRawSize,
			&resultSize,
			&rpcWarningsSize,
			&rpcErrorsSize,
			&errSize,
		)
		if rc != 0 {
			return nil, scrapligoerrors.NewFfiError("poll operation failed", nil)
		}

		if done {
			break
		}

		curPollDelay = getPollDelay(
			curPollDelay,
			n.minPollDelay,
			n.maxPollDelay,
			n.backoffFactor,
		)
	}

	var resultStartTime, resultEndTime uint64

	input := make([]byte, inputSize)

	resultRaw := make([]byte, resultRawSize)

	result := make([]byte, resultSize)

	rpcWarnings := make([]byte, rpcWarningsSize)

	rpcErrors := make([]byte, rpcErrorsSize)

	errString := make([]byte, errSize)

	rc := n.ffiMap.Netconf.FetchOperation(
		n.ptr,
		operationID,
		&resultStartTime,
		&resultEndTime,
		&input,
		&resultRaw,
		&result,
		&rpcWarnings,
		&rpcErrors,
		&errString,
	)
	if rc != 0 {
		return nil, scrapligoerrors.NewFfiError("fetch operation result failed", nil)
	}

	if errSize != 0 {
		return nil, scrapligoerrors.NewFfiError(string(errString), nil)
	}

	return NewResult(
		string(input),
		n.host,
		*n.options.Port,
		resultStartTime,
		resultEndTime,
		resultRaw,
		string(result),
		rpcWarnings,
		rpcErrors,
	), nil
}
