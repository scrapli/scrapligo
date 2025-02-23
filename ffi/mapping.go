package ffi

// Mapping holds mappings to the libscrapli external functions.
type Mapping struct {
	// AssertNoLeaks returns "true" if no leaks were found, otherwise false.
	AssertNoLeaks func() bool
	Driver        DriverMapping
	DriverNetconf DriverNetconfMapping
}

// DriverMapping holds libscrapli mappings specifically for telnet/ssh drivers.
type DriverMapping struct {
	// AllocFromYaml allocates a driver object in zig from the given platform definition yaml file,
	// using the given params.
	AllocFromYaml func(
		platformDefinitionFile string,
		platformVariant string,
		loggerCallback uintptr,
		host string,
		transportKind string,
		port uint16,
		username string,
		password string,
		sessionTimeoutNs uint64,
	) (driverPtr uintptr)
	// AllocFromYamlString ist he same as AllocFromYaml but accepts a loaded YAML string rather than
	// a path to a string. This is useful because a platform definition may be in memory in a go
	// asset.
	AllocFromYamlString func(
		platformDefinitionString string,
		platformVariant string,
		loggerCallback uintptr,
		host string,
		transportKind string,
		port uint16,
		username string,
		password string,
		sessionTimeoutNs uint64,
	) (driverPtr uintptr)
	// Free releases the memory of the driver object at driverPtr -- this should be called *after*
	// Close where possible.
	Free func(driverPtr uintptr)

	// Open opens the telnet/ssh connection of the Driver at driverPtr.
	Open func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) int
	Close func(driverPtr uintptr)

	// PollOperation checks to see if the given operationID is complete -- the state (done or not
	// done) is set into the done bool pointer. If the state is done, the other pointers are also
	// populated such that the Driver "knows" how much space to allocate for the result(raw) and
	// fail/error indicators. Note that while there is a "waitOperation" method in zig, we do *not*
	// use that here as that would block the goroutine -- we simply poll repeatedly until the
	// operation result is ready.
	PollOperation func(
		driverPtr uintptr,
		operationID uint32,
		done *bool,
		resultRawSize,
		resultSize,
		resultFailedIndicatorSize,
		errSize *uint64,
	) int
	// FetchOperation gets the result of the given operationID -- before calling this you must have
	// already understood what the result sizes are such that those pointers can be appropriately
	// allocated for zig to write the results into.
	FetchOperation func(
		driverPtr uintptr,
		operationID uint32,
		resultStartTime *uint64,
		resultEndTime *uint64,
		resultRaw,
		result,
		resultFailedIndicator,
		err *[]byte,
	) int

	// EnterMode submits an EnterMode operation to the underlying driver with the given mode and the
	// configured options. The driver populates the operationID into the uint32 pointer.
	EnterMode func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		requestedMode string,
	) int
	// GetPrompt submits a GetPrompt operation to the underlying driver with the configured options.
	// The driver populates the operationID into the uint32 pointer.
	GetPrompt func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) int
	// SendInput submits a SendInput operation to the underlying driver with the given input and
	// configured options. The driver populates the operationID into the uint32 pointer.
	SendInput func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		input string,
	) int
	// SendPromptedInput submits a SendPromptedInput operation to the underlying driver with the
	// given input, prompt, response, and configured options. The driver populates the operationID
	// into the uint32 pointer.
	SendPromptedInput func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		input string,
		prompt string,
		response string,
		hiddenResponse bool,
		abortInput string,
	) int
}

// DriverNetconfMapping holds libscrapli mappings specifically for the netconf driver.
type DriverNetconfMapping struct {
	// Alloc allocates the driver. See DriverMapping.Alloc for details.
	Alloc func(
		loggerCallback uintptr,
		host string,
		transportKind string,
		port uint16,
		username string,
		password string,
		sessionTimeoutNs uint64,
	) (driverPtr uintptr)
	// Free frees the driver. See DriverMapping.Free for details.
	Free func(driverPtr uintptr)

	// Open opens the driver. See DriverMapping.Open for details.
	Open func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) int
	// Close closes the driver. See DriverMapping.Close for details.
	Close func(driverPtr uintptr)

	// PollOperation polls the given operationID. See DriverMapping.PollerOperation for details.
	PollOperation func(
		driverPtr uintptr,
		operationID uint32,
		done *bool,
		resultRawSize,
		resultSize,
		errSize *uint64,
	) int
	// FetchOperation polls the given operationID. See DriverMapping.FetchOperation for details.
	FetchOperation func(
		driverPtr uintptr,
		operationID uint32,
		resultStartTime *uint64,
		resultEndTime *uint64,
		resultRaw,
		result,
		err *[]byte,
	) int

	// GetConfig submits a GetConfig operation to the underlying driver. The driver populates the
	// operationID into the uint32 pointer.
	GetConfig func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) int
}
