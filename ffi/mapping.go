package ffi

// Mapping holds mappings to the libscrapli external functions.
type Mapping struct {
	// AssertNoLeaks returns "true" if no leaks were found, otherwise false.
	AssertNoLeaks func() bool

	Shared  SharedMapping
	Session SessionMapping
	Cli     CliMapping
	Netconf NetconfMapping
}

// SharedMapping holds common mappings for both cli and netconf drivers.
type SharedMapping struct {
	GetPollFd func(driverPtr uintptr) uint32

	// Free releases the memory of the driver object at driverPtr -- this should be called *after*
	// Close where possible.
	Free func(driverPtr uintptr)

	AllocDriverOptions func() uintptr
	FreeDriverOptions  func(p uintptr)
}

// SessionMapping holds session specific mappings.
type SessionMapping struct {
	Read func(
		driverPtr uintptr,
		buf *[]byte,
		readSize *uint64,
	) uint8
	Write func(
		driverPtr uintptr,
		buf string,
		redacted bool,
	) uint8
	WriteAndReturn func(
		driverPtr uintptr,
		buf string,
		redacted bool,
	) uint8
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
		errSize *uint64,
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

	// EnterMode submits an EnterMode operation to the underlying driver with the given mode and the
	// configured options. The driver populates the operationID into the uint32 pointer.
	EnterMode func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		requestedMode string,
	) uint8
	// GetPrompt submits a GetPrompt operation to the underlying driver with the configured options.
	// The driver populates the operationID into the uint32 pointer.
	GetPrompt func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) uint8
	// SendInput submits a SendInput operation to the underlying driver with the given input and
	// configured options. The driver populates the operationID into the uint32 pointer.
	SendInput func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		input string,
		requestedMode string,
		inputHandling string,
		retainInput bool,
		retainTrailingPrompt bool,
	) uint8
	// SendPromptedInput submits a SendPromptedInput operation to the underlying driver with the
	// given input, prompt, response, and configured options. The driver populates the operationID
	// into the uint32 pointer.
	SendPromptedInput func(
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

// NetconfMapping holds libscrapli mappings specifically for the netconf driver.
type NetconfMapping struct {
	// Alloc allocates the driver. See CliMapping.Alloc for details.
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
		force bool,
	) uint8

	// PollOperation polls the given operationID. See DriverMapping.PollerOperation for details.
	FetchOperationSizes func(
		driverPtr uintptr,
		operationID uint32,
		inputSize,
		resultRawSize,
		resultSize,
		rpcWarningsSize,
		rpcErrorsSize,
		errSize *uint64,
	) uint8
	// FetchOperation polls the given operationID. See CliMapping.FetchOperation for details.
	FetchOperation func(
		driverPtr uintptr,
		operationID uint32,
		resultStartTime *uint64,
		resultEndTime *uint64,
		input,
		resultRaw,
		result,
		rpcWarnings,
		rpcErrors,
		err *[]byte,
	) uint8

	// GetSessionID returns the session id of the current driver session object.
	GetSessionID func(
		driverPtr uintptr,
		sessionID *uint64,
	) uint8

	// GetSubscriptionID writes the subscription id of the given message to the pointer.
	GetSubscriptionID func(
		message string,
		subscriptionID *uint64,
	) uint8

	// GetNextNotificationSize writes the size of the next (if any) notification message into the
	// given size pointer.
	GetNextNotificationSize func(
		driverPtr uintptr,
		size *uint64,
	)

	// GetNextNotificationSize writes the content of the next (if any) notification message into the
	// given message pointer.
	GetNextNotification func(
		driverPtr uintptr,
		notification *[]byte,
	) uint8

	// GetNextSubscriptionSize writes the size of the next (if any) subscription message for the
	// given id into the given size pointer.
	GetNextSubscriptionSize func(
		driverPtr uintptr,
		subscriptionID uint64,
		size *uint64,
	)

	// GetNextSubscription writes the content of the next (if any) subscription message for the
	// given id into the given message pointer.
	GetNextSubscription func(
		driverPtr uintptr,
		subscriptionID uint64,
		subscription *[]byte,
	) uint8

	// RawRPC submits a user defined rpc -- the library will ensure it is properly delimited but the
	// given payload must be valid/correct.
	RawRPC func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		payload string,
		baseNamespacePrefix string,
		extraNamespaces string,
	) uint8

	// GetConfig submits a GetConfig operation to the underlying driver. The driver populates the
	// operationID into the uint32 pointer.
	GetConfig func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		source string,
		filter string,
		filterType string,
		filterNamespacePrefix string,
		filterNamespace string,
		defaultsType string,
	) uint8

	// EditConfig submits an EditConfig operation to the underlying driver. The driver populates the
	// operationID into the uint32 pointer.
	EditConfig func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		config string,
		target string,
		defaultOperation string,
		testOption string,
		errorOption string,
	) uint8

	// CopyConfig submits a CopyConfig operation to the underlying driver. The driver populates the
	// operationID into the uint32 pointer.
	CopyConfig func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		target string,
		source string,
	) uint8

	// DeleteConfig submits a DeleteConfig operation to the underlying driver. The driver populates
	// the operationID into the uint32 pointer.
	DeleteConfig func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		target string,
	) int

	// Lock submits a Lock operation to the underlying driver. The driver populates the operationID
	// into the uint32 pointer.
	Lock func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		target string,
	) uint8

	// Unlock submits an Unlock operation to the underlying driver. The driver populates the
	// operationID into the uint32 pointer.
	Unlock func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		target string,
	) uint8

	// Get submits a Get operation to the underlying driver. The driver populates the operationID
	// into the uint32 pointer.
	Get func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		filter string,
		filterType string,
		filterNamespacePrefix string,
		filterNamespace string,
		defaultsType string,
	) uint8

	// CloseSession submits a CloseSession operation to the underlying driver. The driver populates
	// the operationID into the uint32 pointer.
	CloseSession func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) uint8

	// KillSession submits a KillSession operation to the underlying driver. The driver populates
	// the operationID into the uint32 pointer.
	KillSession func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		sessionID uint64,
	) uint8

	Commit func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) uint8
	Discard func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) uint8
	CancelCommit func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		persistId string,
	) uint8
	Validate func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		source string,
	) uint8

	GetSchema func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		identifier string,
		version string,
		format string,
	) uint8
	GetData func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		datastore,
		filter,
		filterType,
		filterNamespacePrefix,
		filterNamespace,
		configFilter,
		originFilters string,
		maxDepth uint32,
		withOrigin bool,
		defaultsType string,
	) uint8
	EditData func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		datastore string,
		content string,
		defaultOperation string,
	) uint8
	Action func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		action string,
	) uint8
}
