package ffi

// Mapping holds mappings to the libscrapli external functions.
type Mapping struct {
	// AssertNoLeaks returns "true" if no leaks were found, otherwise false.
	AssertNoLeaks func() bool

	Shared  SharedMapping
	Cli     CliMapping
	Netconf NetconfMapping
	Options OptionMapping
}

// SharedMapping holds common mappings for both cli and netconf drivers.
type SharedMapping struct {
	// Free releases the memory of the driver object at driverPtr -- this should be called *after*
	// Close where possible.
	Free func(driverPtr uintptr)

	// Open opens the driver connection of the driver at driverPtr.
	Open func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) int
	Close func(driverPtr uintptr)

	Read func(
		driverPtr uintptr,
		buf *[]byte,
		readSize *uint64,
	)
	Write func(
		driverPtr uintptr,
		buf string,
		redacted bool,
	)
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

	// PollOperation checks to see if the given operationID is complete -- the state (done or not
	// done) is set into the done bool pointer. If the state is done, the other pointers are also
	// populated such that the Cli "knows" how much space to allocate for the result(raw) and
	// fail/error indicators. Note that while there is a "waitOperation" method in zig, we do *not*
	// use that here as that would block the goroutine -- we simply poll repeatedly until the
	// operation result is ready.
	PollOperation func(
		driverPtr uintptr,
		operationID uint32,
		done *bool,
		inputSize,
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
		input,
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
		requestedMode string,
		inputHandling string,
		retainInput bool,
		retainTrailingPrompt bool,
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
		promptPattern string,
		response string,
		hiddenInput bool,
		abortInput string,
		requestedMode string,
		inputHandling string,
		retainTrailingPrompt bool,
	) int
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

	// PollOperation polls the given operationID. See DriverMapping.PollerOperation for details.
	PollOperation func(
		driverPtr uintptr,
		operationID uint32,
		done *bool,
		inputSize,
		resultRawSize,
		resultSize,
		rpcWarningsSize,
		rpcErrorsSize,
		errSize *uint64,
	) int
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
	) int

	// RawRPC submits a user defined rpc -- the library will ensure it is properly delimited but the
	// given payload must be valid/correct.
	RawRPC func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		payload string,
	) int

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
	) int

	// EditConfig submits an EditConfig operation to the underlying driver. The driver populates the
	// operationID into the uint32 pointer.
	EditConfig func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		config string,
		target string,
	) int

	// CopyConfig submits a CopyConfig operation to the underlying driver. The driver populates the
	// operationID into the uint32 pointer.
	CopyConfig func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		target string,
		source string,
	) int

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
	) int

	// Unlock submits an Unlock operation to the underlying driver. The driver populates the
	// operationID into the uint32 pointer.
	Unlock func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		target string,
	) int

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
	) int

	// CloseSession submits a CloseSession operation to the underlying driver. The driver populates
	// the operationID into the uint32 pointer.
	CloseSession func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) int

	// KillSession submits a KillSession operation to the underlying driver. The driver populates
	// the operationID into the uint32 pointer.
	KillSession func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		sessionID uint64,
	) int

	Commit func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) int
	Discard func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) int
	CancelCommit func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
	) int
	Validate func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		source string,
	) int

	GetSchema func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		identifier string,
		version string,
		format string,
	) int
	GetData func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		configFilter string,
		maxDepth uint32,
		withOrigin bool,
		datastore,
		filter,
		filterType,
		filterNamespacePrefix,
		filterNamespace,
		originFilters,
		defaultsType string,
	) int
	EditData func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		datastore string,
		content string,
	) int
	Action func(
		driverPtr uintptr,
		operationID *uint32,
		cancel *bool,
		action string,
	) int
}

// OptionMapping holds libscrapli mappings for the applying driver options.
type OptionMapping struct {
	Session       SessionOptions
	Auth          AuthOptions
	TransportBin  TransportBinOptions
	TransportSSH2 TransportSSH2Options
	TransportTest TransportTestOptions
}

// SessionOptions holds options setters for session related things.
type SessionOptions struct {
	// SetReadSize sets the session read size for the driver at driverPtr.
	SetReadSize func(
		driverPtr uintptr,
		value uint64,
	) int

	// SetReadDelayMinNs sets the session minimum read delay in nanoseconds for the driver
	// at driverPtr.
	SetReadDelayMinNs func(
		driverPtr uintptr,
		value uint64,
	) int

	// SetReadDelayMaxNs sets the session minimum read delay in nanoseconds for the driver
	// at driverPtr.
	SetReadDelayMaxNs func(
		driverPtr uintptr,
		value uint64,
	) int

	// SetReadDelayBackoffFactor sets the backoff factor for the read delay for the driver
	// at driverPtr.
	SetReadDelayBackoffFactor func(
		driverPtr uintptr,
		value uint8,
	) int

	// SetReturnChar sets the return char string for the driver at driverPtr.
	SetReturnChar func(
		driverPtr uintptr,
		value string,
	) int

	// SetOperationTimeoutNs sets the session operation timeout in nanoseconds for the driver
	// at driverPtr.
	SetOperationTimeoutNs func(
		driverPtr uintptr,
		value uint64,
	) int

	// SetOperationMaxSearchDepth sets the maximum search depth to look backward for prompt
	// matching for the driver at driverPtr.
	SetOperationMaxSearchDepth func(
		driverPtr uintptr,
		value uint64,
	) int

	// SetRecorderPath sets the recorder path for the driver at driverPtr.
	SetRecorderPath func(
		driverPtr uintptr,
		value string,
	) int
}

// AuthOptions holds options setters related to authentication.
type AuthOptions struct {
	// SetUsername sets the username for the driver at driverPtr.
	SetUsername func(
		driverPtr uintptr,
		value string,
	) int

	// SetPassword sets the username for the driver at driverPtr.
	SetPassword func(
		driverPtr uintptr,
		value string,
	) int

	// SetPrivateKeyPath sets the private key path for the driver at driverPtr.
	SetPrivateKeyPath func(
		driverPtr uintptr,
		value string,
	) int

	// SetPrivateKeyPassphrase sets the private key passphrase for the driver at driverPtr.
	SetPrivateKeyPassphrase func(
		driverPtr uintptr,
		value string,
	) int

	// SetDriverOptionAuthLookupKeyValue sets a k/v pair in the lookup map for the driver at
	// driverPtr.
	SetDriverOptionAuthLookupKeyValue func(
		driverPtr uintptr,
		key string,
		value string,
	) int

	// SetInSessionAuthBypass sets the in session auth bypass flag for the driver at driverPtr.
	SetInSessionAuthBypass func(
		driverPtr uintptr,
	) int

	// SetUsernamePattern sets the username pcre2 regex pattern for the driver at driverPtr.
	SetUsernamePattern func(
		driverPtr uintptr,
		value string,
	) int

	// SetPasswordPattern sets the password pcre2 regex pattern for the driver at driverPtr.
	SetPasswordPattern func(
		driverPtr uintptr,
		value string,
	) int

	// SetPassphrasePattern sets the passphrase pcre2 regex pattern for the driver at driverPtr.
	SetPassphrasePattern func(
		driverPtr uintptr,
		value string,
	) int
}

// TransportBinOptions holds options setters for the bin transport.
type TransportBinOptions struct {
	// SetBin sets the path to the binary to use when opening the transport for the driver at
	// driverPtr.
	SetBin func(
		driverPtr uintptr,
		value string,
	) int

	// SetExtraOpenArgs sets the extra args to pass when opening the transport for the driver at
	// driverPtr.
	SetExtraOpenArgs func(
		driverPtr uintptr,
		value string,
	) int

	// SetOverrideOpenArgs sets the extra args to pass when opening the transport for the driver at
	// driverPtr.
	SetOverrideOpenArgs func(
		driverPtr uintptr,
		value string,
	) int

	// SetSSHConfigPath sets the ssh config file path for the transport for the driver at driverPtr.
	SetSSHConfigPath func(
		driverPtr uintptr,
		value string,
	) int

	// SetKnownHostsPath sets the ssh config file path for the transport for the driver at
	// driverPtr.
	SetKnownHostsPath func(
		driverPtr uintptr,
		value string,
	) int

	// SetEnableStrictKey sets the flag to enable strict key checking for the driver at driverPtr.
	SetEnableStrictKey func(
		driverPtr uintptr,
	) int

	// SetTermHeight sets the pty term height for the driver at driverPtr.
	SetTermHeight func(
		driverPtr uintptr,
		value uint16,
	) int

	// SetTermWidth sets the pty term width for the driver at driverPtr.
	SetTermWidth func(
		driverPtr uintptr,
		value uint16,
	) int
}

// TransportSSH2Options holds options setters for the ssh2 transport.
type TransportSSH2Options struct {
	// SetLibSSH2Trace enables libssh2 trace for the driver at driverPtr.
	SetLibSSH2Trace func(
		driverPtr uintptr,
	) int
}

// TransportTestOptions holds options setters for the test transport.
type TransportTestOptions struct {
	// SetF sets the "f" (source file) option for the driver at driverPtr.
	SetF func(
		driverPtr uintptr,
		value string,
	) int
}
