package ffi

import "github.com/ebitengine/purego"

func registerCli(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.Cli.Alloc, libScrapliFfi, "ls_cli_alloc")

	purego.RegisterLibFunc(&m.Cli.open, libScrapliFfi, "ls_cli_open")
	purego.RegisterLibFunc(&m.Cli.close, libScrapliFfi, "ls_cli_close")

	purego.RegisterLibFunc(
		&m.Cli.fetchOperationSizes,
		libScrapliFfi,
		"ls_cli_fetch_operation_sizes",
	)
	purego.RegisterLibFunc(&m.Cli.fetchOperation, libScrapliFfi, "ls_cli_fetch_operation")

	purego.RegisterLibFunc(&m.Cli.enterMode, libScrapliFfi, "ls_cli_enter_mode")
	purego.RegisterLibFunc(&m.Cli.getPrompt, libScrapliFfi, "ls_cli_get_prompt")
	purego.RegisterLibFunc(&m.Cli.sendInput, libScrapliFfi, "ls_cli_send_input")
	purego.RegisterLibFunc(&m.Cli.sendInputs, libScrapliFfi, "ls_cli_send_inputs")
	purego.RegisterLibFunc(&m.Cli.sendPromptedInput, libScrapliFfi, "ls_cli_send_prompted_input")

	purego.RegisterLibFunc(&m.Cli.readAny, libScrapliFfi, "ls_cli_read_any")

	purego.RegisterLibFunc(
		&m.Cli.readCallbackShouldExecute,
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

	open func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) uint8

	close func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) uint8

	fetchOperationSizes func(
		driverPtr uintptr,
		operationID uint32,
		operationCount *uint32,
		inputsSize,
		resultsRawSize,
		resultsSize,
		resultsFailedIndicatorSize,
		errSize *uintptr,
	) uint8

	fetchOperation func(
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

	readAny func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) uint8

	readCallbackShouldExecute func(
		buf string,
		name string,
		contains string,
		containsPattern string,
		notContains string,
		execute *bool,
	) uint8
}

// Open opens the driver connection of the driver at driverPtr.
func (m *CliMapping) Open(
	driverPtr uintptr,
	operationID *uint32,
	cancel *bool,
) error {
	return newLibScrapliResult(
		m.open(
			driverPtr,
			operationID,
			cancel,
		),
		"failed to submit open operation",
	).check()
}

// Close closes the cli driver.
func (m *CliMapping) Close(
	driverPtr uintptr,
	operationID *uint32,
	cancel *bool,
) error {
	return newLibScrapliResult(
		m.close(
			driverPtr,
			operationID,
			cancel,
		),
		"failed to submit close operation",
	).check()
}

// FetchOperationSizes gets the result *sizes* for the given operation id.
func (m *CliMapping) FetchOperationSizes(
	driverPtr uintptr,
	operationID uint32,
	operationCount *uint32,
	inputsSize,
	resultsRawSize,
	resultsSize,
	resultsFailedIndicatorSize,
	errSize *uintptr,
) error {
	return newLibScrapliResult(
		m.fetchOperationSizes(
			driverPtr,
			operationID,
			operationCount,
			inputsSize,
			resultsRawSize,
			resultsSize,
			resultsFailedIndicatorSize,
			errSize,
		),
		"fetch operation sizes failed",
	).check()
}

// FetchOperation gets the result of the given operationID -- before calling this you must have
// already understood what the result sizes are such that those pointers can be appropriately
// allocated for zig to write the results into.
func (m *CliMapping) FetchOperation(
	driverPtr uintptr,
	operationID uint32,
	resultStartTime *uint64,
	splits *[]uint64,
	inputs,
	resultsRaw,
	results,
	resultsFailedIndicator,
	err *[]byte,
) error {
	return newLibScrapliResult(
		m.fetchOperation(
			driverPtr,
			operationID,
			resultStartTime,
			splits,
			inputs,
			resultsRaw,
			results,
			resultsFailedIndicator,
			err,
		),
		"fetch operation failed",
	).check()
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
	).check()
}

// ReadAny submit a ReadAny operaiton to the driver.
func (m *CliMapping) ReadAny(
	driverPtr uintptr,
	operationID *uint32,
	cancel *bool,
) error {
	return newLibScrapliResult(
		m.readAny(
			driverPtr,
			operationID,
			cancel,
		),
		"failed to submit readAny operation",
	).check()
}

// ReadCallbackShouldExecute returns true if a readcallback should be executed otherwise false.
func (m *CliMapping) ReadCallbackShouldExecute(
	buf string,
	name string,
	contains string,
	containsPattern string,
	notContains string,
	execute *bool,
) error {
	return newLibScrapliResult(
		m.readCallbackShouldExecute(
			buf,
			name,
			contains,
			containsPattern,
			notContains,
			execute,
		),
		"failed checking if callback should execute",
	).check()
}
