package ffi

import "github.com/ebitengine/purego"

func registerNetconf(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.Netconf.Alloc, libScrapliFfi, "ls_netconf_alloc")

	purego.RegisterLibFunc(&m.Netconf.Open, libScrapliFfi, "ls_netconf_open")
	purego.RegisterLibFunc(&m.Netconf.Close, libScrapliFfi, "ls_netconf_close")

	purego.RegisterLibFunc(
		&m.Netconf.FetchOperationSizes,
		libScrapliFfi,
		"ls_netconf_fetch_operation_sizes",
	)
	purego.RegisterLibFunc(
		&m.Netconf.FetchOperation,
		libScrapliFfi,
		"ls_netconf_fetch_operation",
	)

	purego.RegisterLibFunc(&m.Netconf.GetSessionID, libScrapliFfi, "ls_netconf_get_session_id")
	purego.RegisterLibFunc(
		&m.Netconf.GetSubscriptionID,
		libScrapliFfi,
		"ls_netconf_get_subscription_id",
	)

	purego.RegisterLibFunc(
		&m.Netconf.GetNextNotificationSize,
		libScrapliFfi,
		"ls_netconf_next_notification_message_size",
	)
	purego.RegisterLibFunc(
		&m.Netconf.GetNextNotification,
		libScrapliFfi,
		"ls_netconf_next_notification_message",
	)

	purego.RegisterLibFunc(
		&m.Netconf.GetNextSubscriptionSize,
		libScrapliFfi,
		"ls_netconf_next_subscription_message_size",
	)
	purego.RegisterLibFunc(
		&m.Netconf.GetNextSubscription,
		libScrapliFfi,
		"ls_netconf_next_subscription_message",
	)

	purego.RegisterLibFunc(&m.Netconf.RawRPC, libScrapliFfi, "ls_netconf_raw_rpc")

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
		errSize *uintptr,
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
