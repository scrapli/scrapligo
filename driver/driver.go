package driver

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

	d.options.DefinitionString = string(definitionBytes)

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

	minNs := scrapligoconstants.DefaultReadDelayMinNs
	maxNs := scrapligoconstants.DefaultReadDelayMaxNs

	d.minPollDelay = minNs * scrapligoconstants.ReadDelayMultiplier
	if d.options.Session.ReadDelayMinNs != nil {
		d.minPollDelay = *d.options.Session.ReadDelayMinNs * scrapligoconstants.ReadDelayMultiplier
	}

	d.maxPollDelay = maxNs * scrapligoconstants.ReadDelayMultiplier
	if d.options.Session.ReadDelayMaxNs != nil {
		d.maxPollDelay = *d.options.Session.ReadDelayMaxNs * scrapligoconstants.ReadDelayMultiplier
	}

	d.backoffFactor = scrapligoconstants.DefaultReadDelayBackoffFactor
	if d.options.Session.ReadDelayBackoffFactor != nil {
		d.backoffFactor = *d.options.Session.ReadDelayBackoffFactor
	}

	return d, nil
}

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

		d.ffiMap.Driver.Free(d.ptr)
	}()

	var loggerCallback uintptr
	if d.options.LoggerCallback != nil {
		loggerCallback = purego.NewCallback(d.options.LoggerCallback)
	}

	d.ptr = d.ffiMap.Driver.Alloc(
		d.options.DefinitionString,
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

	status := d.ffiMap.Driver.Open(d.ptr, &operationID, &cancel)
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
func (d *Driver) Close() {
	if d.ptr == 0 {
		return
	}

	d.ffiMap.Driver.Close(d.ptr)
	d.ffiMap.Driver.Free(d.ptr)
}

func getPollDelay(curVal, minVal, maxVal uint64, backoffFactor uint8) uint64 {
	newVal := curVal
	newVal *= uint64(backoffFactor)

	if newVal > maxVal {
		newVal = maxVal
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

	var resultRawSize, resultSize, resultFailedIndicatorSize, errSize uint64

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

		rc := d.ffiMap.Driver.PollOperation(
			d.ptr,
			operationID,
			&done,
			&resultRawSize,
			&resultSize,
			&resultFailedIndicatorSize,
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

	resultRaw := make([]byte, resultRawSize)

	result := make([]byte, resultSize)

	resultFailedWhenIndicator := make([]byte, resultFailedIndicatorSize)

	errString := make([]byte, errSize)

	rc := d.ffiMap.Driver.FetchOperation(
		d.ptr,
		operationID,
		&resultStartTime,
		&resultEndTime,
		&resultRaw,
		&result,
		&resultFailedWhenIndicator,
		&errString,
	)
	if rc != 0 {
		return nil, scrapligoerrors.NewFfiError("fetch operation result failed", nil)
	}

	if errSize != 0 {
		return nil, scrapligoerrors.NewFfiError(string(errString), nil)
	}

	return NewResult(
		"",
		d.host,
		*d.options.Port,
		resultStartTime,
		resultEndTime,
		resultRaw,
		string(result),
		resultFailedWhenIndicator,
	), nil
}
