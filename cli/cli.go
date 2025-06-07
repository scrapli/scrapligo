package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	scrapligoassets "github.com/scrapli/scrapligo/assets"
	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
	scrapligointernal "github.com/scrapli/scrapligo/internal"
	scrapligologging "github.com/scrapli/scrapligo/logging"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	"golang.org/x/sys/unix"
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

// Cli is an object representing a connection to a device of some sort -- this object wraps the
// underlying zig driver (created via libscrapli).
type Cli struct {
	ptr     uintptr
	pollFd  int
	ffiMap  *scrapligoffi.Mapping
	host    string
	options *scrapligointernal.Options
}

// NewCli returns a new instance of Cli setup with the given options. The definitionFileOrName
// should be the name of one of the platforms that has a definition embedded in this package's
// assets, or a file path to a valid yaml definition.
func NewCli[T PlatformNameOrString](
	definitionFileOrName T,
	host string,
	opts ...scrapligooptions.Option,
) (*Cli, error) {
	ffiMap, err := scrapligoffi.GetMapping()
	if err != nil {
		return nil, err
	}

	c := &Cli{
		ffiMap:  ffiMap,
		host:    host,
		options: scrapligointernal.NewOptions(),
	}

	definitionBytes, err := getDefinitionBytes(definitionFileOrName)
	if err != nil {
		return nil, err
	}

	c.options.Driver.DefinitionString = string(definitionBytes)

	for _, opt := range opts {
		err = opt(c.options)
		if err != nil {
			return nil, scrapligoerrors.NewOptionsError("failed applying option", err)
		}
	}

	if c.options.Port == nil {
		var p uint16

		switch c.options.TransportKind { //nolint: exhaustive
		case scrapligointernal.TransportKindTelnet:
			p = scrapligoconstants.DefaultTelnetPort
		default:
			p = scrapligoconstants.DefaultSSHPort
		}

		c.options.Port = &p
	}

	return c, nil
}

// GetPtr returns the pointer to the zig driver, don't use this unless you know what you are doing,
// this is just exposed so you *can* get to it if you want to.
func (c *Cli) GetPtr() (uintptr, *scrapligoffi.Mapping) {
	return c.ptr, c.ffiMap
}

// Open opens the driver object. This method spawns the underlying zig driver which the Cli then
// holds a pointer to. All Cli operations operate against this pointer (though this is
// transparent to the user).
func (c *Cli) Open(ctx context.Context) (*Result, error) {
	// ensure we dealloc if something happens, otherwise users calls to defer close would not be
	// super handy
	cleanup := false

	defer func() {
		if !cleanup {
			return
		}

		c.ffiMap.Shared.Free(c.ptr)
	}()

	c.ptr = c.ffiMap.Cli.Alloc(
		c.options.Driver.DefinitionString,
		scrapligologging.LoggerToLoggerCallback(
			c.options.Logger,
			uint8(scrapligologging.IntFromLevel(c.options.LoggerLevel)),
		),
		c.host,
		*c.options.Port,
		string(c.options.TransportKind),
	)

	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("failed to allocate cli", nil)
	}

	c.pollFd = int(c.ffiMap.Shared.GetPollFd(c.ptr))
	if c.pollFd == 0 {
		return nil, scrapligoerrors.NewFfiError("failed to allocate cli", nil)
	}

	err := c.options.Apply(c.ptr, c.ffiMap)
	if err != nil {
		return nil, scrapligoerrors.NewFfiError("failed to applying cli options", err)
	}

	cancel := false

	var operationID uint32

	status := c.ffiMap.Cli.Open(c.ptr, &operationID, &cancel)
	if status != 0 {
		cleanup = true

		return nil, scrapligoerrors.NewFfiError("failed to submit open operation", nil)
	}

	result, err := c.getResult(ctx, &cancel, operationID)
	if err != nil {
		cleanup = true

		return nil, err
	}

	return result, nil
}

// Close closes the driver object. This also deallocates the underlying (zig) driver object.
func (c *Cli) Close(ctx context.Context) (*Result, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	// as long as driver ptr is not 0 we *always* want to free
	defer c.ffiMap.Shared.Free(c.ptr)

	cancel := false

	var operationID uint32

	status := c.ffiMap.Cli.Close(c.ptr, &operationID, &cancel)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit close operation", nil)
	}

	return c.getResult(ctx, &cancel, operationID)
}

func (c *Cli) getResult(
	ctx context.Context,
	cancel *bool,
	operationID uint32,
) (*Result, error) {
	done := make(chan struct{}, 1)
	defer close(done)

	var operationCount uint32

	// so in go flavor we actually use ctx to cause libscrapli to timeout vs python where we rely on
	// the timeouts in libscrapli itself. so in this case we need to ensure that we do not block the
	// context so it can properly cancel on timeout/cancellation...
	go func() {
		select {
		case <-ctx.Done():
			*cancel = true

			return
		case <-done:
			return
		}
	}()

	pollFd := &unix.FdSet{}
	pollFd.Set(c.pollFd)

	var n int

	for {
		var err error

		n, err = unix.Select(c.pollFd+1, pollFd, &unix.FdSet{}, &unix.FdSet{}, nil)
		if err != nil {
			if errors.Is(err, unix.EINTR) {
				// python automagically handles interrupts i guess go doesnt, so just act like
				// we do on the python side when polling the wakeup fd
				continue
			}

			return nil, scrapligoerrors.NewFfiError("waiting on operation ready signal", err)
		}

		break
	}

	// if the context wasn't cancelled the goroutine will still be running, this will stop it
	done <- struct{}{}

	out := make([]byte, n)

	_, _ = unix.Read(c.pollFd, out)

	var inputsSize, resultsRawSize, resultsSize, resultsFailedIndicatorSize, errSize uint64

	rc := c.ffiMap.Cli.FetchOperationSizes(
		c.ptr,
		operationID,
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

	var resultStartTime uint64

	splits := make([]uint64, operationCount)

	inputs := make([]byte, inputsSize)

	resultsRaw := make([]byte, resultsRawSize)

	results := make([]byte, resultsSize)

	resultsFailedWhenIndicator := make([]byte, resultsFailedIndicatorSize)

	errString := make([]byte, errSize)

	rc = c.ffiMap.Cli.FetchOperation(
		c.ptr,
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
		c.host,
		*c.options.Port,
		inputs,
		resultStartTime,
		splits,
		resultsRaw,
		results,
		resultsFailedWhenIndicator,
	), nil
}
