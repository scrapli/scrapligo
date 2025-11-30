package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	scrapligoassets "github.com/scrapli/scrapligo/assets"
	scrapligoclidefinitionoptions "github.com/scrapli/scrapligo/cli/definitionoptions"
	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
	scrapligointernal "github.com/scrapli/scrapligo/internal"
	scrapligologging "github.com/scrapli/scrapligo/logging"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	"golang.org/x/sys/unix"
)

type loadedDefinition struct {
	content      []byte
	platformName string
}

func getDefinitionBytes(definitionFileOrName string) (*loadedDefinition, error) {
	d := &loadedDefinition{}

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
			d.content, err = scrapligoassets.Assets.ReadFile(
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

			d.platformName = platformName
		}
	}

	if len(d.content) == 0 {
		// didn't load from assets, so we'll try to load the file
		d.content, err = os.ReadFile(definitionFileOrNameString) //nolint: gosec
		if err != nil {
			return nil, scrapligoerrors.NewUtilError(
				fmt.Sprintf("failed loading definition file at path %q", definitionFileOrName),
				err,
			)
		}

		d.platformName = strings.TrimSuffix(
			filepath.Base(definitionFileOrNameString),
			filepath.Ext(definitionFileOrNameString),
		)
	}

	return d, nil
}

// Cli is an object representing a connection to a device of some sort -- this object wraps the
// underlying zig driver (created via libscrapli).
type Cli struct {
	ptr     uintptr
	pollFd  int
	ffiMap  *scrapligoffi.Mapping
	host    string
	options *scrapligointernal.Options
	l       *scrapligologging.AnyLogger
}

// NewCli returns a new instance of Cli setup with the given options.
func NewCli(
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

	for _, opt := range opts {
		err = opt(c.options)
		if err != nil {
			return nil, scrapligoerrors.NewOptionsError("failed applying option", err)
		}
	}

	c.l = c.options.GetLogger()

	loadedDef, err := getDefinitionBytes(c.options.Cli.DefinitionFileOrName)
	if err != nil {
		return nil, err
	}

	c.options.Cli.DefinitionString = string(loadedDef.content)

	if c.options.Port == 0 {
		var p uint16

		switch c.options.TransportKind { //nolint: exhaustive
		case scrapligointernal.TransportKindTelnet:
			p = scrapligoconstants.DefaultTelnetPort
		default:
			p = scrapligoconstants.DefaultSSHPort
		}

		c.options.Port = p
	}

	if !c.options.Cli.SkipStaticOptions {
		// for platforms that have... quirks, its difficult to fully encapsulate setting up a
		// connection in purely yaml... so... there are py/go "extensions" in the
		// scrapli_definitions project that are pulled into scrapli/scrapligo in order to facilitate
		// these quirks -- this includes options, things like mikrotik that *really* wants you to
		// modify a username with some extra chars to change how the device behaves, here is where
		// we apply those options. obviously this can be skipped with the appropriate option.
		for _, opt := range scrapligoclidefinitionoptions.GetPlatformOptions().
			GetOptionsForPlatform(loadedDef.platformName) {
			err = opt(c.options)
			if err != nil {
				return nil, scrapligoerrors.NewOptionsError("failed applying (static) option", err)
			}
		}
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

	optionsPtr := c.ffiMap.Shared.AllocDriverOptions()
	defer c.ffiMap.Shared.FreeDriverOptions(optionsPtr)

	c.options.Apply(optionsPtr)

	c.ptr = c.ffiMap.Cli.Alloc(
		c.host,
		optionsPtr,
	)

	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("failed to allocate cli", nil)
	}

	c.pollFd = int(c.ffiMap.Shared.GetPollFd(c.ptr))
	if c.pollFd == 0 {
		return nil, scrapligoerrors.NewFfiError("failed to allocate cli", nil)
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

// Write writes the input to the session -- this bypasses the driver operation loop in zig, use
// with caution.
func (c *Cli) Write(input string) error {
	if c.ptr == 0 {
		return scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	status := c.ffiMap.Session.Write(c.ptr, input, false)
	if status != 0 {
		return scrapligoerrors.NewFfiError("failed executing write", nil)
	}

	return nil
}

// WriteAndReturn writes the given input and then sends a return character -- this bypasses the
// driver operation loop in zig, use with caution.
func (c *Cli) WriteAndReturn(input string) error {
	if c.ptr == 0 {
		return scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	status := c.ffiMap.Session.WriteAndReturn(c.ptr, input, false)
	if status != 0 {
		return scrapligoerrors.NewFfiError("failed executing writeAndReturn", nil)
	}

	return nil
}

// WriteReturn writes a return character -- this bypasses the driver operation loop in zig, use
// with caution.
func (c *Cli) WriteReturn() error {
	if c.ptr == 0 {
		return scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	status := c.ffiMap.Session.WriteReturn(c.ptr)
	if status != 0 {
		return scrapligoerrors.NewFfiError("failed executing writeReturn", nil)
	}

	return nil
}

// Read reads from the session -- this bypasses the driver operation loop in zig, use with caution.
func (c *Cli) Read() ([]byte, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	var buf []byte

	var readSize uint64

	status := c.ffiMap.Session.Read(c.ptr, &buf, &readSize)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed executing read", nil)
	}

	return buf[0:readSize], nil
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
		// always wrap the context error (even if nil) so we catch cancels/deadline exceeded and
		// users can errors.Is with that
		return nil, scrapligoerrors.NewFfiError(string(errString), ctx.Err())
	}

	return NewResult(
		c.host,
		c.options.Port,
		inputs,
		resultStartTime,
		splits,
		resultsRaw,
		results,
		resultsFailedWhenIndicator,
	), nil
}
