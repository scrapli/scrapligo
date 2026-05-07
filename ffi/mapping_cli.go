package ffi

import "github.com/ebitengine/purego"

func registerCli(m *Mapping, libScrapliFfi uintptr) {
	// ENHANCEMENT?: is it possible to have my own register funcs that bypass reflection?
	//  driver creation/destruction
	purego.RegisterLibFunc(&m.Cli.Alloc, libScrapliFfi, "ls_cli_alloc")

	purego.RegisterLibFunc(&m.Cli.Open, libScrapliFfi, "ls_cli_open")
	purego.RegisterLibFunc(&m.Cli.Close, libScrapliFfi, "ls_cli_close")

	purego.RegisterLibFunc(
		&m.Cli.FetchOperationSizes,
		libScrapliFfi,
		"ls_cli_fetch_operation_sizes",
	)
	purego.RegisterLibFunc(&m.Cli.FetchOperation, libScrapliFfi, "ls_cli_fetch_operation")

	purego.RegisterLibFunc(&m.Cli.enterMode, libScrapliFfi, "ls_cli_enter_mode")
	purego.RegisterLibFunc(&m.Cli.getPrompt, libScrapliFfi, "ls_cli_get_prompt")
	purego.RegisterLibFunc(&m.Cli.sendInput, libScrapliFfi, "ls_cli_send_input")
	purego.RegisterLibFunc(&m.Cli.sendInputs, libScrapliFfi, "ls_cli_send_inputs")
	purego.RegisterLibFunc(&m.Cli.sendPromptedInput, libScrapliFfi, "ls_cli_send_prompted_input")

	purego.RegisterLibFunc(&m.Cli.ReadAny, libScrapliFfi, "ls_cli_read_any")
	purego.RegisterLibFunc(
		&m.Cli.ReadCallbackShouldExecute,
		libScrapliFfi,
		"ls_cli_read_callback_should_execute",
	)
}

// CliMapping holds libscrapli mappings specifically for cli drivers.
type CliMapping struct {
	// Alloc allocates a driver object in zig.
	Alloc func(
		host string,
		optionsPtr uintptr,
	) (driverPtr uintptr)

	// Open opens the driver connection of the driver at driverPtr.
	Open func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) uint8
	Close func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) uint8

	FetchOperationSizes func(
		driverPtr uintptr,
		operationID uint32,
		operationCount *uint32,
		inputsSize,
		resultsRawSize,
		resultsSize,
		resultsFailedIndicatorSize,
		errSize *uintptr,
	) uint8
	// FetchOperation gets the result of the given operationID -- before calling this you must have
	// already understood what the result sizes are such that those pointers can be appropriately
	// allocated for zig to write the results into.
	FetchOperation func(
		driverPtr uintptr,
		operationID uint32,
		resultStartTime *uint64,
		splits *[]uint64,
		inputs,
		resultsRaw,
		results,
		resultsFailedIndicator,
		err *[]byte,
	) uint8

	enterMode func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		requestedMode string,
	) uint8

	getPrompt func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) uint8

	sendInput func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		input string,
		requestedMode string,
		inputHandling string,
		retainInput bool,
		retainTrailingPrompt bool,
	) uint8

	sendInputs func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		inputs string,
		requestedMode string,
		inputHandling string,
		retainInput bool,
		retainTrailingPrompt bool,
	) uint8

	sendPromptedInput func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		input string,
		prompt string,
		promptPattern string,
		response string,
		abortInput string,
		requestedMode string,
		inputHandling string,
		hiddenInput bool,
		retainTrailingPrompt bool,
	) uint8

	// ReadAny submit a ReadAny operaiton to the driver.
	ReadAny func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) uint8

	// ReadCallbackShouldExecute returns true if a readcallback should be executed otherwise false.
	ReadCallbackShouldExecute func(
		buf string,
		name string,
		contains string,
		containsPattern string,
		notContains string,
		execute *bool,
	) uint8
}

// EnterMode submits an EnterMode operation to the underlying driver with the given mode and the
// configured options. The driver populates the operationID into the uint32 pointer.
func (m *CliMapping) EnterMode(
	driverPtr uintptr,
	operationID *uint32,
	cancel *bool,
	requestedMode string,
) error {
	return newLibScrapliResult(
		m.enterMode(
			driverPtr,
			operationID,
			cancel,
			requestedMode,
		),
		"failed to submit enterMode operation",
		nil,
	).check()
}

// GetPrompt submits a GetPrompt operation to the underlying driver with the configured options.
// The driver populates the operationID into the uint32 pointer.
func (m *CliMapping) GetPrompt(
	driverPtr uintptr,
	operationID *uint32,
	cancel *bool,
) error {
	return newLibScrapliResult(
		m.getPrompt(
			driverPtr,
			operationID,
			cancel,
		),
		"failed to submit getPrompt operation",
		nil,
	).check()
}

// SendInput submits a SendInput operation to the underlying driver with the given input and
// configured options. The driver populates the operationID into the uint32 pointer.
func (m *CliMapping) SendInput(
	driverPtr uintptr,
	operationID *uint32,
	cancel *bool,
	input string,
	requestedMode string,
	inputHandling string,
	retainInput bool,
	retainTrailingPrompt bool,
) error {
	return newLibScrapliResult(
		m.sendInput(
			driverPtr,
			operationID,
			cancel,
			input,
			requestedMode,
			inputHandling,
			retainInput,
			retainTrailingPrompt,
		),
		"failed to submit sendInput operation",
		nil,
	).check()
}

// SendInputs submits a SendInputs operation to the underlying driver with the given input and
// configured options. The driver populates the operationID into the uint32 pointer.
func (m *CliMapping) SendInputs(
	driverPtr uintptr,
	operationID *uint32,
	cancel *bool,
	inputs string,
	requestedMode string,
	inputHandling string,
	retainInput bool,
	retainTrailingPrompt bool,
) error {
	return newLibScrapliResult(
		m.sendInputs(
			driverPtr,
			operationID,
			cancel,
			inputs,
			requestedMode,
			inputHandling,
			retainInput,
			retainTrailingPrompt,
		),
		"failed to submit sendInputs operation",
		nil,
	).check()
}

// SendPromptedInput submits a SendPromptedInput operation to the underlying driver with the
// given input, prompt, response, and configured options. The driver populates the operationID
// into the uint32 pointer.
func (m *CliMapping) SendPromptedInput(
	driverPtr uintptr,
	operationID *uint32,
	cancel *bool,
	input string,
	prompt string,
	promptPattern string,
	response string,
	abortInput string,
	requestedMode string,
	inputHandling string,
	hiddenInput bool,
	retainTrailingPrompt bool,
) error {
	return newLibScrapliResult(
		m.sendPromptedInput(
			driverPtr,
			operationID,
			cancel,
			input,
			prompt,
			promptPattern,
			response,
			abortInput,
			requestedMode,
			inputHandling,
			hiddenInput,
			retainTrailingPrompt,
		),
		"failed to submit sendPromptedInput operation",
		nil,
	).check()
}
