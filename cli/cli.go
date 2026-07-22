package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	scrapligoassets "github.com/scrapli/scrapligo/v2/assets"
	scrapligoclidefinitionoptions "github.com/scrapli/scrapligo/v2/cli/definitionoptions"
	scrapligoconstants "github.com/scrapli/scrapligo/v2/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
	scrapligoffi "github.com/scrapli/scrapligo/v2/ffi"
	scrapligointernal "github.com/scrapli/scrapligo/v2/internal"
	scrapligologging "github.com/scrapli/scrapligo/v2/logging"
	scrapligooptions "github.com/scrapli/scrapligo/v2/options"
	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
	"golang.org/x/sys/unix"
)

const defaultCancelCleanupGrace = 5 * time.Second

var errOperationCancelCleanupTimeout = errors.New("libscrapli operation cleanup timeout")

func loadDefinition(o *scrapligointernal.Options) error {
	definitionFileOrNameString := o.Cli.DefinitionFileOrName

	var err error

	assetPlatformNames := GetPlatformNames()

	for _, platformName := range assetPlatformNames {
		if platformName != definitionFileOrNameString {
			continue
		}

		b, err := scrapligoassets.Assets.ReadFile(
			fmt.Sprintf("definitions/%s.yaml", platformName),
		)
		if err != nil {
			return scrapligoerrors.NewUtilError(
				fmt.Sprintf(
					"failed loading definition asset for platform %q",
					o.Cli.DefinitionFileOrName,
				),
				err,
			)
		}

		o.Cli.DefinitionPlatform = platformName
		o.Cli.DefinitionString = string(b)

		return nil
	}

	// didn't load from assets, so we'll try to load the file
	b, err := os.ReadFile(definitionFileOrNameString) //nolint: gosec
	if err != nil {
		return scrapligoerrors.NewUtilError(
			fmt.Sprintf("failed loading definition file at path %q", o.Cli.DefinitionFileOrName),
			err,
		)
	}

	o.Cli.DefinitionPlatform = strings.TrimSuffix(
		filepath.Base(definitionFileOrNameString),
		filepath.Ext(definitionFileOrNameString),
	)
	o.Cli.DefinitionString = string(b)

	return nil
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

	if c.options.Cli.DefinitionPlatform == "" {
		err := loadDefinition(c.options)
		if err != nil {
			return nil, err
		}
	}

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
			GetOptionsForPlatform(c.options.Cli.DefinitionPlatform) {
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

// GetOptions returns the options as supplied in a json-ish string -- should only really be used
// for testing as it doesnt really serve any purpose otherwise and costs some allocations and calls
// across the ffi boundary. Exposed for testing reasons.
func (c *Cli) GetOptions() (string, error) {
	optionsPtr := c.ffiMap.Shared.AllocDriverOptions()
	defer c.ffiMap.Shared.FreeDriverOptions(optionsPtr)

	c.options.Apply(optionsPtr)

	var optionsSize uintptr

	err := c.ffiMap.Shared.FetchOptionsSize(
		optionsPtr,
		&optionsSize,
	)
	if err != nil {
		return "", err
	}

	optionsStr := make([]byte, optionsSize)

	err = c.ffiMap.Shared.FetchOptions(
		optionsPtr,
		&optionsStr,
	)
	if err != nil {
		return "", err
	}

	return string(optionsStr), nil
}

// Open opens the driver object. This method spawns the underlying zig driver which the Cli then
// holds a pointer to. All Cli operations operate against this pointer (though this is
// transparent to the user).
func (c *Cli) Open(ctx context.Context) (*Result, error) {
	// ensure we dealloc if something happens, otherwise users calls to defer close would not be
	// super handy
	cleanup, cleanupFree := false, true

	defer func() {
		if !cleanup {
			return
		}

		if cleanupFree {
			c.ffiMap.Shared.Free(c.ptr)
		}

		c.ptr = 0
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
		cleanup = true

		return nil, scrapligoerrors.NewFfiError("failed to allocate cli", nil)
	}

	cancel := false

	var operationID uint32

	err := c.ffiMap.Cli.Open(c.ptr, &operationID, &cancel)
	if err != nil {
		cleanup = true

		return nil, err
	}

	result, err := c.getResult(ctx, &cancel, operationID)
	if err != nil {
		cleanup = true
		cleanupFree = !errors.Is(err, errOperationCancelCleanupTimeout)

		return nil, err
	}

	return result, nil
}

// Close closes the driver object. This also deallocates the underlying (zig) driver object.
func (c *Cli) Close(ctx context.Context) (*Result, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cleanup := true

	defer func() {
		if cleanup {
			c.ffiMap.Shared.Free(c.ptr)
		}

		c.ptr = 0
	}()

	cancel := false

	var operationID uint32

	err := c.ffiMap.Cli.Close(c.ptr, &operationID, &cancel)
	if err != nil {
		return nil, err
	}

	result, err := c.getResult(ctx, &cancel, operationID)
	if err != nil {
		cleanup = !errors.Is(err, errOperationCancelCleanupTimeout)

		return nil, err
	}

	return result, nil
}

// ReplaceDefinition replaces the "definition" of the driver. Most importantly changes/updates
// the prompt pattern, but also updates the modes etc. available in the driver.
func (c *Cli) ReplaceDefinition(definitionFileOrString string) error {
	if c.ptr == 0 {
		return scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	err := loadDefinition(c.options)
	if err != nil {
		return err
	}

	c.options.Cli.DefinitionFileOrName = definitionFileOrString

	return c.ffiMap.Cli.ReplaceDefinition(c.ptr, c.options.Cli.DefinitionString)
}

func (c *Cli) getResult( //nolint: funlen,gocyclo
	ctx context.Context,
	cancel *bool,
	operationID uint32,
) (*Result, error) {
	var operationCount uint32

	var n int
	var ctxErr error
	var cancelCleanupDeadline time.Time

	pollFds := []unix.PollFd{{Fd: int32(c.pollFd), Events: unix.POLLIN}} //nolint: gosec

	for {
		if ctxErr == nil && ctx.Err() != nil {
			ctxErr = ctx.Err()
			*cancel = true

			cancelCleanupDeadline = time.Now().Add(c.cancelCleanupGrace())
		}

		if ctxErr != nil && time.Now().After(cancelCleanupDeadline) {
			msg := "libscrapli operation did not finish after cancellation; " +
				"skipping free to avoid blocking in FFI cleanup"

			return nil, scrapligoerrors.NewFfiError(
				msg,
				errors.Join(ctxErr, errOperationCancelCleanupTimeout),
			)
		}

		pollFds[0].Revents = 0

		var err error

		n, err = unix.Poll(pollFds, scrapligoconstants.ReadyFDPollTimeoutMs)
		if err != nil {
			if errors.Is(err, unix.EINTR) {
				// python automagically handles interrupts i guess go doesnt, so just act like
				// we do on the python side when polling the wakeup fd
				continue
			}

			return nil, scrapligoerrors.NewFfiError("waiting on operation ready signal", err)
		}

		if n > 0 {
			if pollFds[0].Revents&unix.POLLNVAL != 0 {
				return nil, scrapligoerrors.NewFfiError(
					"waiting on operation ready signal",
					unix.EBADF,
				)
			}

			break
		}
	}

	var out [1]byte

	for {
		_, err := unix.Read(c.pollFd, out[:])
		if err == nil {
			break
		}

		if errors.Is(err, unix.EINTR) {
			// same as loop above -- retry on interrupts
			continue
		}

		return nil, scrapligoerrors.NewFfiError("draining operation ready signal", err)
	}

	var (
		inputsSize                 uintptr
		resultsRawSize             uintptr
		resultsSize                uintptr
		resultsFailedIndicatorSize uintptr
		errSize                    uintptr
		lastErrStrSize             uintptr
	)

	err := c.ffiMap.Cli.FetchOperationSizes(
		c.ptr,
		operationID,
		&operationCount,
		&inputsSize,
		&resultsRawSize,
		&resultsSize,
		&resultsFailedIndicatorSize,
		&errSize,
		&lastErrStrSize,
	)
	if err != nil {
		return nil, err
	}

	var resultStartTime uint64

	splits := make([]uint64, operationCount)

	inputs := make([]byte, inputsSize)

	resultsRaw := make([]byte, resultsRawSize)

	results := make([]byte, resultsSize)

	resultsFailedWhenIndicator := make([]byte, resultsFailedIndicatorSize)

	errString := make([]byte, errSize)

	lastErrString := make([]byte, lastErrStrSize)

	err = c.ffiMap.Cli.FetchOperation(
		c.ptr,
		operationID,
		&resultStartTime,
		&splits,
		&inputs,
		&resultsRaw,
		&results,
		&resultsFailedWhenIndicator,
		&errString,
		&lastErrString,
	)
	if err != nil {
		return nil, err
	}

	if errSize != 0 {
		// always wrap the context error (even if nil) so we catch cancels/deadline exceeded and
		// users can errors.Is with that
		outErrMsg := string(errString)

		if lastErrStrSize > 0 {
			outErrMsg += fmt.Sprintf(": %s", string(lastErrString))
		}

		return nil, scrapligoerrors.NewFfiError(outErrMsg, ctx.Err())
	}

	if ctxErr != nil {
		return nil, ctxErr
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

func (c *Cli) cancelCleanupGrace() time.Duration {
	if c.options == nil || c.options.Session.OperationTimeoutNs == nil {
		return defaultCancelCleanupGrace
	}

	v := *c.options.Session.OperationTimeoutNs
	if v == 0 {
		return defaultCancelCleanupGrace
	}

	return time.Duration(scrapligoutil.SafeUint64ToInt64(v))
}
