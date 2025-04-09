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

// NewDriver returns a new instance of Driver setup with the given options.
func NewDriver(
	host string,
	opts ...scrapligooptions.Option,
) (*Driver, error) {
	ffiMap, err := scrapligoffi.GetMapping()
	if err != nil {
		return nil, err
	}

	n := &Driver{
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

	minNs := scrapligoconstants.DefaultReadDelayMinNs
	maxNs := scrapligoconstants.DefaultReadDelayMaxNs

	n.minPollDelay = minNs * scrapligoconstants.ReadDelayMultiplier
	if n.options.Session.ReadDelayMinNs != nil {
		n.minPollDelay = *n.options.Session.ReadDelayMinNs * scrapligoconstants.ReadDelayMultiplier
	}

	n.maxPollDelay = maxNs * scrapligoconstants.ReadDelayMultiplier
	if n.options.Session.ReadDelayMaxNs != nil {
		n.maxPollDelay = *n.options.Session.ReadDelayMaxNs * scrapligoconstants.ReadDelayMultiplier
	}

	n.backoffFactor = scrapligoconstants.DefaultReadDelayBackoffFactor
	if n.options.Session.ReadDelayBackoffFactor != nil {
		n.backoffFactor = *n.options.Session.ReadDelayBackoffFactor
	}

	return n, nil
}

// Driver is an object representing a netconf connection to a device of some sort -- this object
// wraps the underlying zig (netconf) driver (created via libscrapli).
type Driver struct {
	ptr     uintptr
	ffiMap  *scrapligoffi.Mapping
	host    string
	options *scrapligointernal.Options

	minPollDelay  uint64
	maxPollDelay  uint64
	backoffFactor uint8
}

// GetPtr returns the pointer to the zig driver, don't use this unless you know what you are doing,
// this is just exposed so you *can* get to it if you want to.
func (d *Driver) GetPtr() (uintptr, *scrapligoffi.Mapping) {
	return d.ptr, d.ffiMap
}

// Open opens the driver object. This method spawns the underlying zig driver which the Cli then
// holds a pointer to. All Cli operations operate against this pointer (though this is
// transparent to the user).
func (d *Driver) Open(ctx context.Context) (*Result, error) {
	// ensure we dealloc if something happens, otherwise users calls to defer close would not be
	// super handy
	cleanup := false

	defer func() {
		if !cleanup {
			return
		}

		d.ffiMap.Shared.Free(d.ptr)
	}()

	var loggerCallback uintptr
	if d.options.LoggerCallback != nil {
		loggerCallback = purego.NewCallback(d.options.LoggerCallback)
	}

	d.ptr = d.ffiMap.Netconf.Alloc(
		loggerCallback,
		d.host,
		*d.options.Port,
		string(d.options.TransportKind),
	)

	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("failed to allocate driver", nil)
	}

	err := d.options.Apply(d.ptr, d.ffiMap)
	if err != nil {
		return nil, scrapligoerrors.NewFfiError("failed to applying driver options", err)
	}

	cancel := false

	var operationID uint32

	status := d.ffiMap.Shared.Open(d.ptr, &operationID, &cancel)
	if status != 0 {
		cleanup = true

		return nil, scrapligoerrors.NewFfiError("failed to submit open operation", nil)
	}

	result, err := d.getResult(ctx, &cancel, operationID)
	if err != nil {
		cleanup = true

		return nil, err
	}

	return result, nil
}

// Close closes the netconf object. This also deallocates the underlying (zig) netconf object.
func (d *Driver) Close() {
	if d.ptr == 0 {
		return
	}

	d.ffiMap.Shared.Close(d.ptr)
	d.ffiMap.Shared.Free(d.ptr)
}

// GetSessionID returns the session-id as parsed during the capabilities exchange -- if we for some
// reason didn't parse the session-id during capabilities exchange this will return an error.
func (d *Driver) GetSessionID() (uint64, error) {
	if d.ptr == 0 {
		return 0, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	var sessionID uint64

	status := d.ffiMap.Netconf.GetSessionID(d.ptr, &sessionID)
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

func (d *Driver) getResult(
	ctx context.Context,
	cancel *bool,
	operationID uint32,
) (*Result, error) {
	var done bool

	var inputSize, resultRawSize, resultSize, rpcWarningsSize, rpcErrorsSize, errSize uint64

	minNs := scrapligoconstants.DefaultReadDelayMinNs

	curPollDelay := minNs * scrapligoconstants.ReadDelayMultiplier
	if d.options.Session.ReadDelayMinNs != nil {
		curPollDelay = *d.options.Session.ReadDelayMinNs * scrapligoconstants.ReadDelayMultiplier
	}

	for {
		select {
		case <-ctx.Done():
			*cancel = true

			return nil, ctx.Err()
		default:
		}

		// we obviously cant have too tight a loop here or cpu will go nuts and we'll block things,
		// so we'll sleep the same as the zig read delay will be
		curPollDelay = getPollDelay(
			curPollDelay,
			d.minPollDelay,
			d.maxPollDelay,
			d.backoffFactor,
		)

		time.Sleep(time.Duration(scrapligoutil.SafeUint64ToInt64(curPollDelay)))

		rc := d.ffiMap.Netconf.PollOperation(
			d.ptr,
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

		if !done {
			continue
		}

		break
	}

	var resultStartTime, resultEndTime uint64

	input := make([]byte, inputSize)

	resultRaw := make([]byte, resultRawSize)

	result := make([]byte, resultSize)

	rpcWarnings := make([]byte, rpcWarningsSize)

	rpcErrors := make([]byte, rpcErrorsSize)

	errString := make([]byte, errSize)

	rc := d.ffiMap.Netconf.FetchOperation(
		d.ptr,
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
		d.host,
		*d.options.Port,
		resultStartTime,
		resultEndTime,
		resultRaw,
		string(result),
		rpcWarnings,
		rpcErrors,
	), nil
}
