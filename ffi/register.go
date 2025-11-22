package ffi

import "github.com/ebitengine/purego"

func registerShared(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.Shared.GetPollFd, libScrapliFfi, "ls_shared_get_poll_fd")
	purego.RegisterLibFunc(&m.Shared.Free, libScrapliFfi, "ls_shared_free")

	purego.RegisterLibFunc(&m.Shared.AllocDriverOptions, libScrapliFfi, "ls_alloc_driver_options")
	purego.RegisterLibFunc(&m.Shared.FreeDriverOptions, libScrapliFfi, "ls_free_driver_options")
}

func registerSession(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.Session.Read, libScrapliFfi, "ls_session_read")
	purego.RegisterLibFunc(&m.Session.Write, libScrapliFfi, "ls_session_write")
	purego.RegisterLibFunc(&m.Session.WriteAndReturn, libScrapliFfi, "ls_session_write_and_return")
}

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

	purego.RegisterLibFunc(&m.Cli.EnterMode, libScrapliFfi, "ls_cli_enter_mode")
	purego.RegisterLibFunc(&m.Cli.GetPrompt, libScrapliFfi, "ls_cli_get_prompt")
	purego.RegisterLibFunc(&m.Cli.SendInput, libScrapliFfi, "ls_cli_send_input")
	purego.RegisterLibFunc(&m.Cli.SendPromptedInput, libScrapliFfi, "ls_cli_send_prompted_input")

	purego.RegisterLibFunc(&m.Cli.ReadAny, libScrapliFfi, "ls_cli_read_any")
	purego.RegisterLibFunc(
		&m.Cli.ReadCallbackShouldExecute,
		libScrapliFfi,
		"ls_cli_read_callback_should_execute",
	)
}

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
