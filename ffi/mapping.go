package ffi

// Mapping holds mappings to the libscrapli external functions.
type Mapping struct {
	// AssertNoLeaks returns "true" if no leaks were found, otherwise false.
	AssertNoLeaks func() bool

	Shared  SharedMapping
	Session SessionMapping
	Cli     CliMapping
	Netconf NetconfMapping
	Options OptionMapping
}

// SharedMapping holds common mappings for both cli and netconf drivers.
type SharedMapping struct {
	GetPollFd func(driverPtr uintptr) uint32

	// Free releases the memory of the driver object at driverPtr -- this should be called *after*
	// Close where possible.
	Free func(driverPtr uintptr)
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
	// Alloc allocates a driver object in zig -- it expects *all* the possible options.
	Alloc func(
		definitionString string,
		loggerCallback uintptr,
		host string,
		port uint16,
		transportKind string,
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
		loggerCallback uintptr,
		host string,
		port uint16,
		transportKind string,
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
	) uint8
	Action func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		action string,
	) uint8
}

// OptionMapping holds libscrapli mappings for the applying driver options.
type OptionMapping struct {
	Session       SessionOptions
	Auth          AuthOptions
	TransportBin  TransportBinOptions
	TransportSSH2 TransportSSH2Options
	TransportTest TransportTestOptions
	Netconf       NetconfOptions
}

// SessionOptions holds options setters for session related things.
type SessionOptions struct {
	// SetReadSize sets the session read size for the driver at driverPtr.
	SetReadSize func(
		driverPtr uintptr,
		value uint64,
	) uint8

	// SetReadMinDelayNs sets the session minimum read delay in ns.
	SetReadMinDelayNs func(
		driverPtr uintptr,
		value uint64,
	) uint8

	// SetReadMaxDelayNs sets the session maximum read delay in ns.
	SetReadMaxDelayNs func(
		driverPtr uintptr,
		value uint64,
	) uint8

	// SetReturnChar sets the return char string for the driver at driverPtr.
	SetReturnChar func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetOperationTimeoutNs sets the session operation timeout in nanoseconds for the driver
	// at driverPtr.
	SetOperationTimeoutNs func(
		driverPtr uintptr,
		value uint64,
	) uint8

	// SetOperationMaxSearchDepth sets the maximum search depth to look backward for prompt
	// matching for the driver at driverPtr.
	SetOperationMaxSearchDepth func(
		driverPtr uintptr,
		value uint64,
	) uint8

	// SetRecordDestination sets the record destination path for the driver at driverPtr.
	SetRecordDestination func(
		driverPtr uintptr,
		value string,
	) uint8
}

// AuthOptions holds options setters related to authentication.
type AuthOptions struct {
	// SetUsername sets the username for the driver at driverPtr.
	SetUsername func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetPassword sets the username for the driver at driverPtr.
	SetPassword func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetPrivateKeyPath sets the private key path for the driver at driverPtr.
	SetPrivateKeyPath func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetPrivateKeyPassphrase sets the private key passphrase for the driver at driverPtr.
	SetPrivateKeyPassphrase func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetDriverOptionAuthLookupKeyValue sets a k/v pair in the lookup map for the driver at
	// driverPtr.
	SetDriverOptionAuthLookupKeyValue func(
		driverPtr uintptr,
		key string,
		value string,
	) uint8

	// SetForceInSessionAuth sets the in session auth bypass flag for the driver at driverPtr.
	SetForceInSessionAuth func(
		driverPtr uintptr,
	) uint8

	// SetBypassInSessionAuth sets the in session auth bypass flag for the driver at driverPtr.
	SetBypassInSessionAuth func(
		driverPtr uintptr,
	) uint8

	// SetUsernamePattern sets the username pcre2 regex pattern for the driver at driverPtr.
	SetUsernamePattern func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetPasswordPattern sets the password pcre2 regex pattern for the driver at driverPtr.
	SetPasswordPattern func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetPassphrasePattern sets the passphrase pcre2 regex pattern for the driver at driverPtr.
	SetPassphrasePattern func(
		driverPtr uintptr,
		value string,
	) uint8
}

// TransportBinOptions holds options setters for the bin transport.
type TransportBinOptions struct {
	// SetBin sets the path to the binary to use when opening the transport for the driver at
	// driverPtr.
	SetBin func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetExtraOpenArgs sets the extra args to pass when opening the transport for the driver at
	// driverPtr.
	SetExtraOpenArgs func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetOverrideOpenArgs sets the extra args to pass when opening the transport for the driver at
	// driverPtr.
	SetOverrideOpenArgs func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetSSHConfigPath sets the ssh config file path for the transport for the driver at driverPtr.
	SetSSHConfigPath func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetKnownHostsPath sets the ssh config file path for the transport for the driver at
	// driverPtr.
	SetKnownHostsPath func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetEnableStrictKey sets the flag to enable strict key checking for the driver at driverPtr.
	SetEnableStrictKey func(
		driverPtr uintptr,
	) uint8

	// SetTermHeight sets the pty term height for the driver at driverPtr.
	SetTermHeight func(
		driverPtr uintptr,
		value uint16,
	) uint8

	// SetTermWidth sets the pty term width for the driver at driverPtr.
	SetTermWidth func(
		driverPtr uintptr,
		value uint16,
	) uint8
}

// TransportSSH2Options holds options setters for the ssh2 transport.
type TransportSSH2Options struct {
	// SetKnownHostsPath sets the known hosts path for the driver at driverPtr.
	SetKnownHostsPath func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetLibSSH2Trace enables libssh2 trace for the driver at driverPtr.
	SetLibSSH2Trace func(
		driverPtr uintptr,
	) uint8

	// SetProxyJumpHost sets the proxy jump host for the driver at driverPtr.
	SetProxyJumpHost func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetProxyJumpPort sets the proxy jump port for the driver at driverPtr.
	SetProxyJumpPort func(
		driverPtr uintptr,
		value uint16,
	) uint8

	// SetProxyJumpUsername sets the proxy jump username for the driver at driverPtr.
	SetProxyJumpUsername func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetProxyJumpPassword sets the proxy jump password for the driver at driverPtr.
	SetProxyJumpPassword func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetProxyJumpPrivateKeyPath sets the proxy jump private key path for the driver at driverPtr.
	SetProxyJumpPrivateKeyPath func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetProxyJumpPrivateKeyPassphrase sets the proxy jump private key passhrase for the driver at
	// driverPtr.
	SetProxyJumpPrivateKeyPassphrase func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetProxyJumpLibSSH2Trace enables libssh2 trace for the proxy jump driver at driverPtr.
	SetProxyJumpLibSSH2Trace func(
		driverPtr uintptr,
	) uint8
}

// TransportTestOptions holds options setters for the test transport.
type TransportTestOptions struct {
	// SetF sets the "f" (source file) option for the driver at driverPtr.
	SetF func(
		driverPtr uintptr,
		value string,
	) uint8
}

// NetconfOptions holds options setters for netconf objects.
type NetconfOptions struct {
	// SetErrorTag sets the error tag option for the driver at driverPtr. Default is "rpc-error>".
	SetErrorTag func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetPreferredVersion sets the preferred netconf version for the driver at driverPtr.
	SetPreferredVersion func(
		driverPtr uintptr,
		value string,
	) uint8

	// SetMessagePollIntervalNS sets the message poll interval in ns for the driver at driverPtr.
	SetMessagePollIntervalNS func(
		driverPtr uintptr,
		value uint64,
	) uint8
}
