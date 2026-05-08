package ffi

import "github.com/ebitengine/purego"

func registerNetconf(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.Netconf.Alloc, libScrapliFfi, "ls_netconf_alloc")

	purego.RegisterLibFunc(&m.Netconf.open, libScrapliFfi, "ls_netconf_open")
	purego.RegisterLibFunc(&m.Netconf.close, libScrapliFfi, "ls_netconf_close")

	purego.RegisterLibFunc(
		&m.Netconf.fetchOperationSizes,
		libScrapliFfi,
		"ls_netconf_fetch_operation_sizes",
	)
	purego.RegisterLibFunc(
		&m.Netconf.fetchOperation,
		libScrapliFfi,
		"ls_netconf_fetch_operation",
	)

	purego.RegisterLibFunc(&m.Netconf.getSessionID, libScrapliFfi, "ls_netconf_get_session_id")
	purego.RegisterLibFunc(
		&m.Netconf.getSubscriptionID,
		libScrapliFfi,
		"ls_netconf_get_subscription_id",
	)

	purego.RegisterLibFunc(
		&m.Netconf.getNextNotificationSize,
		libScrapliFfi,
		"ls_netconf_next_notification_message_size",
	)
	purego.RegisterLibFunc(
		&m.Netconf.getNextNotification,
		libScrapliFfi,
		"ls_netconf_next_notification_message",
	)

	purego.RegisterLibFunc(
		&m.Netconf.getNextSubscriptionSize,
		libScrapliFfi,
		"ls_netconf_next_subscription_message_size",
	)
	purego.RegisterLibFunc(
		&m.Netconf.getNextSubscription,
		libScrapliFfi,
		"ls_netconf_next_subscription_message",
	)

	purego.RegisterLibFunc(&m.Netconf.rawRPC, libScrapliFfi, "ls_netconf_raw_rpc")

	purego.RegisterLibFunc(&m.Netconf.GetConfig, libScrapliFfi, "ls_netconf_get_config")
	purego.RegisterLibFunc(&m.Netconf.EditConfig, libScrapliFfi, "ls_netconf_edit_config")
	purego.RegisterLibFunc(&m.Netconf.CopyConfig, libScrapliFfi, "ls_netconf_copy_config")
	purego.RegisterLibFunc(&m.Netconf.DeleteConfig, libScrapliFfi, "ls_netconf_delete_config")
	purego.RegisterLibFunc(&m.Netconf.Lock, libScrapliFfi, "ls_netconf_lock")
	purego.RegisterLibFunc(&m.Netconf.Unlock, libScrapliFfi, "ls_netconf_unlock")
	purego.RegisterLibFunc(&m.Netconf.Get, libScrapliFfi, "ls_netconf_get")
	purego.RegisterLibFunc(&m.Netconf.CloseSession, libScrapliFfi, "ls_netconf_close_session")
	purego.RegisterLibFunc(&m.Netconf.KillSession, libScrapliFfi, "ls_netconf_kill_session")

	purego.RegisterLibFunc(&m.Netconf.Commit, libScrapliFfi, "ls_netconf_commit")
	purego.RegisterLibFunc(&m.Netconf.Discard, libScrapliFfi, "ls_netconf_discard")
	purego.RegisterLibFunc(&m.Netconf.CancelCommit, libScrapliFfi, "ls_netconf_cancel_commit")
	purego.RegisterLibFunc(&m.Netconf.Validate, libScrapliFfi, "ls_netconf_validate")

	purego.RegisterLibFunc(&m.Netconf.GetSchema, libScrapliFfi, "ls_netconf_get_schema")
	purego.RegisterLibFunc(&m.Netconf.GetData, libScrapliFfi, "ls_netconf_get_data")
	purego.RegisterLibFunc(&m.Netconf.EditData, libScrapliFfi, "ls_netconf_edit_data")
	purego.RegisterLibFunc(&m.Netconf.Action, libScrapliFfi, "ls_netconf_action")
}

// NetconfMapping holds libscrapli mappings specifically for the netconf driver.
type NetconfMapping struct {
	// Alloc allocates the driver. See CliMapping.Alloc for details.
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
		force bool,
	) uint8

	fetchOperationSizes func(
		driverPtr uintptr,
		operationID uint32,
		inputSize,
		resultRawSize,
		resultSize,
		rpcWarningsSize,
		rpcErrorsSize,
		errSize *uintptr,
	) uint8

	fetchOperation func(
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

	getSessionID func(
		driverPtr uintptr,
		sessionID *uint64,
	) uint8

	getSubscriptionID func(
		message string,
		subscriptionID *uint64,
	) uint8

	getNextNotificationSize func(
		driverPtr uintptr,
		size *uint64,
	) uint8

	getNextNotification func(
		driverPtr uintptr,
		notification *[]byte,
	) uint8

	getNextSubscriptionSize func(
		driverPtr uintptr,
		subscriptionID uint64,
		size *uint64,
	) uint8

	getNextSubscription func(
		driverPtr uintptr,
		subscriptionID uint64,
		subscription *[]byte,
	) uint8

	rawRPC func(
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
	) uint8

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

// Open opens the driver connection of the driver at driverPtr.
func (m *NetconfMapping) Open(
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
func (m *NetconfMapping) Close(
	driverPtr uintptr,
	operationID *uint32,
	cancel *bool,
	force bool,
) error {
	return newLibScrapliResult(
		m.close(
			driverPtr,
			operationID,
			cancel,
			force,
		),
		"failed to submit close operation",
	).check()
}

// FetchOperationSizes gets the result *sizes* for the given operation id.
func (m *NetconfMapping) FetchOperationSizes(
	driverPtr uintptr,
	operationID uint32,
	inputSize,
	resultRawSize,
	resultSize,
	rpcWarningsSize,
	rpcErrorsSize,
	errSize *uintptr,
) error {
	return newLibScrapliResult(
		m.fetchOperationSizes(
			driverPtr,
			operationID,
			inputSize,
			resultRawSize,
			resultSize,
			rpcWarningsSize,
			rpcErrorsSize,
			errSize,
		),
		"fetch operation sizes failed",
	).check()
}

// FetchOperation gets the result of the given operationID -- before calling this you must have
// already understood what the result sizes are such that those pointers can be appropriately
// allocated for zig to write the results into.
func (m *NetconfMapping) FetchOperation(
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
) error {
	return newLibScrapliResult(
		m.fetchOperation(
			driverPtr,
			operationID,
			resultStartTime,
			resultEndTime,
			input,
			resultRaw,
			result,
			rpcWarnings,
			rpcErrors,
			err,
		),
		"fetch operation failed",
	).check()
}

// GetSessionID returns the session id of the current driver session object.
func (m *NetconfMapping) GetSessionID(
	driverPtr uintptr,
	sessionID *uint64,
) error {
	return newLibScrapliResult(
		m.getSessionID(
			driverPtr,
			sessionID,
		),
		"fetch session-id failed",
	).check()
}

// GetSubscriptionID writes the subscription id of the given message to the pointer.
func (m *NetconfMapping) GetSubscriptionID(
	message string,
	subscriptionID *uint64,
) error {
	return newLibScrapliResult(
		m.getSubscriptionID(
			message,
			subscriptionID,
		),
		"fetch subscription-id failed",
	).check()
}

// GetNextNotificationSize writes the size of the next (if any) notification message into the
// given size pointer.
func (m *NetconfMapping) GetNextNotificationSize(
	driverPtr uintptr,
	size *uint64,
) error {
	return newLibScrapliResult(
		m.getNextNotificationSize(
			driverPtr,
			size,
		),
		"fetch next notification size failed",
	).check()
}

// GetNextNotification writes the content of the next (if any) notification message into the
// given message pointer.
func (m *NetconfMapping) GetNextNotification(
	driverPtr uintptr,
	notification *[]byte,
) error {
	return newLibScrapliResult(
		m.getNextNotification(
			driverPtr,
			notification,
		),
		"fetch next notification failed",
	).check()
}

// GetNextSubscriptionSize writes the size of the next (if any) subscription message for the
// given id into the given size pointer.
func (m *NetconfMapping) GetNextSubscriptionSize(
	driverPtr uintptr,
	subscriptionID uint64,
	size *uint64,
) error {
	return newLibScrapliResult(
		m.getNextSubscriptionSize(
			driverPtr,
			subscriptionID,
			size,
		),
		"fetch subscription size failed",
	).check()
}

// GetNextSubscription writes the content of the next (if any) subscription message for the
// given id into the given message pointer.
func (m *NetconfMapping) GetNextSubscription(
	driverPtr uintptr,
	subscriptionID uint64,
	subscription *[]byte,
) error {
	return newLibScrapliResult(
		m.getNextSubscription(
			driverPtr,
			subscriptionID,
			subscription,
		),
		"fetch subscription failed",
	).check()
}

// RawRPC submits a user defined rpc -- the library will ensure it is properly delimited but the
// given payload must be valid/correct.
func (m *NetconfMapping) RawRPC(
	driverPtr uintptr,
	operationID *uint32,
	cancel *bool,
	payload string,
	baseNamespacePrefix string,
	extraNamespaces string,
) error {
	return newLibScrapliResult(
		m.rawRPC(
			driverPtr,
			operationID,
			cancel,
			payload,
			baseNamespacePrefix,
			extraNamespaces,
		),
		"failed to submit raw rpc operation",
	).check()
}
