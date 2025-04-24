package cli

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ebitengine/purego"
	scrapligoassets "github.com/scrapli/scrapligo/assets"
	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
	scrapligointernal "github.com/scrapli/scrapligo/internal"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligoutil "github.com/scrapli/scrapligo/util"
)

// PlatformNameOrString is a string-like interface so you can pass a PlatformName or "normal" string
// to the driver constructor.
type PlatformNameOrString interface {
	~string
}

func getDefinitionBytes[T PlatformNameOrString](definitionFileOrName T) ([]byte, error) {
	var definitionBytes []byte

	var definitionFileOrNameString string

	var err error

	switch v := any(definitionFileOrName).(type) {
	case PlatformName:
		definitionFileOrNameString = v.String()
	case string:
		definitionFileOrNameString = v
	}

	assetPlatformNames := GetPlatformNames()

	for _, platformName := range assetPlatformNames {
		if platformName == definitionFileOrNameString {
			definitionBytes, err = scrapligoassets.Assets.ReadFile(
				fmt.Sprintf("definitions/%s.yaml", platformName),
			)
			if err != nil {
				return nil, scrapligoerrors.NewUtilError(
					fmt.Sprintf(
						"failed loading definition asset for platform %q",
						definitionFileOrName,
					),
					err,
				)
			}
		}
	}

	if len(definitionBytes) == 0 {
		// didn't load from assets, so we'll try to load the file
		definitionBytes, err = os.ReadFile(definitionFileOrNameString) //nolint: gosec
		if err != nil {
			return nil, scrapligoerrors.NewUtilError(
				fmt.Sprintf("failed loading definition file at path %q", definitionFileOrName),
				err,
			)
		}
	}

	return definitionBytes, nil
}

// TODO rename to *NewCli*
// NewDriver returns a new instance of Driver setup with the given options. The definitionFileOrName
// should be the name of one of the platforms that has a definition embedded in this package's
// assets, or a file path to a valid yaml definition.
func NewDriver[T PlatformNameOrString](
	definitionFileOrName T,
	host string,
	opts ...scrapligooptions.Option,
) (*Driver, error) {
	ffiMap, err := scrapligoffi.GetMapping()
	if err != nil {
		return nil, err
	}

	d := &Driver{
		ffiMap:  ffiMap,
		host:    host,
		options: scrapligointernal.NewOptions(),
	}

	definitionBytes, err := getDefinitionBytes(definitionFileOrName)
	if err != nil {
		return nil, err
	}

	d.options.Driver.DefinitionString = string(definitionBytes)

	for _, opt := range opts {
		err = opt(d.options)
		if err != nil {
			return nil, scrapligoerrors.NewOptionsError("failed applying option", err)
		}
	}

	if d.options.Port == nil {
		var p uint16

		switch d.options.TransportKind { //nolint: exhaustive
		case scrapligointernal.TransportKindTelnet:
			p = scrapligoconstants.DefaultTelnetPort
		default:
			p = scrapligoconstants.DefaultSSHPort
		}

		d.options.Port = &p
	}

	d.minPollDelay = scrapligoconstants.DefaultReadDelayMinNs
	if d.options.Session.ReadDelayMinNs != nil {
		d.minPollDelay = *d.options.Session.ReadDelayMinNs
	}

	d.maxPollDelay = scrapligoconstants.DefaultReadDelayMaxNs
	if d.options.Session.ReadDelayMaxNs != nil {
		d.maxPollDelay = *d.options.Session.ReadDelayMaxNs
	}

	d.backoffFactor = scrapligoconstants.DefaultReadDelayBackoffFactor
	if d.options.Session.ReadDelayBackoffFactor != nil {
		d.backoffFactor = *d.options.Session.ReadDelayBackoffFactor
	}

	return d, nil
}

// TODO rename to Cli
// Driver is an object representing a connection to a device of some sort -- this object wraps the
// underlying zig driver (created via libscrapli).
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

// Open opens the driver object. This method spawns the underlying zig driver which the Driver then
// holds a pointer to. All Driver operations operate against this pointer (though this is
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

	d.ptr = d.ffiMap.Cli.Alloc(
		d.options.Driver.DefinitionString,
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

// Close closes the driver object. This also deallocates the underlying (zig) driver object.
func (d *Driver) Close(ctx context.Context) (*Result, error) {
	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	// as long as driver ptr is not 0 we *always* want to free
	defer d.ffiMap.Shared.Free(d.ptr)

	cancel := false

	var operationID uint32

	status := d.ffiMap.Shared.Close(d.ptr, &operationID, &cancel)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit close operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}

func getPollDelay(curVal, minVal, maxVal uint64, backoffFactor uint8) uint64 {
	newVal := curVal
	newVal *= uint64(backoffFactor)

	if newVal > maxVal {
		// we backoff up to max then reset when we are over
		newVal = minVal
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

	var operationCount uint32

	var inputsSize, resultsRawSize, resultsSize, resultsFailedIndicatorSize, errSize uint64

	// start w/ max delay, its literally 0.001s so its not much and almost no chance anything is
	// already ready to fetch in that time anyway! also, once/if we hit max dealy we reset to min
	// delay before backing off back to max. so this causes us to have one slow check then fast
	// checks up to the max again
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

		rc := d.ffiMap.Cli.PollOperation(
			d.ptr,
			operationID,
			&done,
			&operationCount,
			&inputsSize,
			&resultsRawSize,
			&resultsSize,
			&resultsFailedIndicatorSize,
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
			d.minPollDelay,
			d.maxPollDelay,
			d.backoffFactor,
		)
	}

	var resultStartTime uint64

	splits := make([]uint64, operationCount)

	inputs := make([]byte, inputsSize)

	resultsRaw := make([]byte, resultsRawSize)

	results := make([]byte, resultsSize)

	resultsFailedWhenIndicator := make([]byte, resultsFailedIndicatorSize)

	errString := make([]byte, errSize)

	rc := d.ffiMap.Cli.FetchOperation(
		d.ptr,
		operationID,
		&resultStartTime,
		&splits,
		&inputs,
		&resultsRaw,
		&results,
		&resultsFailedWhenIndicator,
		&errString,
	)
	if rc != 0 {
		return nil, scrapligoerrors.NewFfiError("fetch operation result failed", nil)
	}

	if errSize != 0 {
		return nil, scrapligoerrors.NewFfiError(string(errString), nil)
	}

	return NewResult(
		d.host,
		*d.options.Port,
		inputs,
		resultStartTime,
		splits,
		resultsRaw,
		results,
		resultsFailedWhenIndicator,
	), nil
}
