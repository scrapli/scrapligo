package driver

import (
	"context"
	"time"

	"github.com/ebitengine/purego"
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
)

// NewDriver returns a new instance of Driver setup with the given options.
func NewDriver(
	platformDefinitionFile,
	host string,
	opts ...Option,
) (*Driver, error) {
	ffiMap, err := scrapligoffi.GetMapping()
	if err != nil {
		return nil, err
	}

	d := &Driver{
		ffiMap: ffiMap,

		platformDefinitionFile: platformDefinitionFile,

		host: host,

		options: newOptions(),
	}

	for _, opt := range opts {
		err = opt(d)
		if err != nil {
			return nil, scrapligoerrors.NewOptionsError("failed applying option", err)
		}
	}

	if d.options.port == nil {
		var p uint16

		switch d.options.transportKind { //nolint: exhaustive
		case TransportKindTelnet:
			p = DefaultTelnetPort
		default:
			p = DefaultSSHPort
		}

		d.options.port = &p
	}

	return d, nil
}

// Driver is an object representing a connection to a device of some sort -- this object wraps the
// underlying zig driver (created via libscrapli).
type Driver struct {
	ptr    uintptr
	ffiMap *scrapligoffi.Mapping

	platformDefinitionFile string

	host string

	options options
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
	if d.options.loggerCallback != nil {
		loggerCallback = purego.NewCallback(d.options.loggerCallback)
	}

	d.ptr = d.ffiMap.Driver.AllocFromYaml(
		d.platformDefinitionFile,
		d.options.platformVariant,
		loggerCallback,
		d.host,
		string(d.options.transportKind),
		*d.options.port,
		d.options.username,
		d.options.password,
		// timeouts governed by contexts in go!
		0,
	)

	if d.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("failed to allocate driver", nil)
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

	// TODO we can run native language callbacks here if we want?

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

func (d *Driver) getResult(ctx context.Context, cancel *bool, operationID uint32) (*Result, error) {
	var done bool

	var resultRawSize, resultSize, resultFailedIndicatorSize, errSize uint64

	for {
		select {
		case <-ctx.Done():
			*cancel = true

			return nil, ctx.Err()
		default:
		}

		// TODO obviously figuring out how tight to loop this will be impactful af on things i think
		// i think it probably makes sense to make it exactly the same as the read delay?
		time.Sleep(time.Millisecond)

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
		*d.options.port,
		resultStartTime,
		resultEndTime,
		resultRaw,
		string(result),
	), nil
}

// EnterMode is used to explicitly enter a mode (i.e. enter "config mode" or "shell" or some other
// platform specific "mode").
func (d *Driver) EnterMode(ctx context.Context, requestedMode string) (*Result, error) {
	cancel := false

	var operationID uint32

	status := d.ffiMap.Driver.EnterMode(d.ptr, &operationID, &cancel, requestedMode)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit enterMode operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}

// GetPrompt returns a Result object containing the current "prompt" of the target device.
func (d *Driver) GetPrompt(ctx context.Context) (*Result, error) {
	cancel := false

	var operationID uint32

	status := d.ffiMap.Driver.GetPrompt(d.ptr, &operationID, &cancel)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit getPrompt operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}

// SendInput sends an "input" to the device. Historically scrapli(go) had "SendCommand(s)" and
// "SendConfig(s)" operations, but these no longer exist. Instead, we have SendInput or SendInputs
// which accept their respective options -- the options can (among other things) control the "mode"
// (historically "privilege level") at which to send the input(s).
func (d *Driver) SendInput(ctx context.Context, input string) (*Result, error) {
	cancel := false

	var operationID uint32

	status := d.ffiMap.Driver.SendInput(d.ptr, &operationID, &cancel, input)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit sendInput operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}

// SendInputs send multiple "inputs" to the device.
func (d *Driver) SendInputs(ctx context.Context, inputs ...string) (*MultiResult, error) {
	cancel := false

	results := NewMultiResult(d.host, *d.options.port)

	for _, input := range inputs {
		var operationID uint32

		status := d.ffiMap.Driver.SendInput(d.ptr, &operationID, &cancel, input)
		if status != 0 {
			return nil, scrapligoerrors.NewFfiError("failed to submit sendInput operation", nil)
		}

		r, err := d.getResult(ctx, &cancel, operationID)
		if err != nil {
			return nil, err
		}

		results.AppendResult(r)
	}

	return results, nil
}

// SendPromptedInput sends an `input` to the device expecting the given `prompt`, finally sending
// the `response`.
func (d *Driver) SendPromptedInput(
	ctx context.Context,
	input,
	prompt,
	response string,
) (*Result, error) {
	cancel := false

	var operationID uint32

	status := d.ffiMap.Driver.SendPromptedInput(
		d.ptr,
		&operationID,
		&cancel,
		input,
		prompt,
		response,
		false,
		"",
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit sendPromptedInput operation", nil)
	}

	return d.getResult(ctx, &cancel, operationID)
}
