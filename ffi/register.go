package ffi

import "github.com/ebitengine/purego"

func registerShared(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.Shared.Free, libScrapliFfi, "ls_shared_free")

	purego.RegisterLibFunc(&m.Shared.Read, libScrapliFfi, "ls_shared_read_session")
	purego.RegisterLibFunc(&m.Shared.Write, libScrapliFfi, "ls_shared_write_session")
}

func registerCli(m *Mapping, libScrapliFfi uintptr) {
	// ENHANCEMENT?: is it possible to have my own register funcs that bypass reflection?
	//  driver creation/destruction
	purego.RegisterLibFunc(&m.Cli.Alloc, libScrapliFfi, "ls_cli_alloc")

	purego.RegisterLibFunc(&m.Cli.Open, libScrapliFfi, "ls_cli_open")
	purego.RegisterLibFunc(&m.Cli.Close, libScrapliFfi, "ls_cli_close")

	purego.RegisterLibFunc(&m.Cli.PollOperation, libScrapliFfi, "ls_cli_poll_operation")
	purego.RegisterLibFunc(&m.Cli.FetchOperation, libScrapliFfi, "ls_cli_fetch_operation")

	purego.RegisterLibFunc(&m.Cli.EnterMode, libScrapliFfi, "ls_cli_enter_mode")
	purego.RegisterLibFunc(&m.Cli.GetPrompt, libScrapliFfi, "ls_cli_get_prompt")
	purego.RegisterLibFunc(&m.Cli.SendInput, libScrapliFfi, "ls_cli_send_input")
	purego.RegisterLibFunc(&m.Cli.SendPromptedInput, libScrapliFfi, "ls_cli_send_prompted_input")
}

func registerNetconf(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.Netconf.Alloc, libScrapliFfi, "ls_netconf_alloc")

	purego.RegisterLibFunc(&m.Netconf.Open, libScrapliFfi, "ls_netconf_open")
	purego.RegisterLibFunc(&m.Netconf.Close, libScrapliFfi, "ls_netconf_close")

	purego.RegisterLibFunc(
		&m.Netconf.PollOperation,
		libScrapliFfi,
		"ls_netconf_poll_operation",
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

func registerOptions(m *Mapping, libScrapliFfi uintptr) { //nolint: funlen
	// session
	purego.RegisterLibFunc(
		&m.Options.Session.SetReadSize,
		libScrapliFfi,
		"ls_option_session_read_size",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetReadDelayMinNs,
		libScrapliFfi,
		"ls_option_session_read_delay_min_ns",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetReadDelayMaxNs,
		libScrapliFfi,
		"ls_option_session_read_delay_max_ns",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetReadDelayBackoffFactor,
		libScrapliFfi,
		"ls_option_session_read_delay_backoff_factor",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetReturnChar,
		libScrapliFfi,
		"ls_option_session_return_char",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetOperationTimeoutNs,
		libScrapliFfi,
		"ls_option_session_operation_timeout_ns",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetOperationMaxSearchDepth,
		libScrapliFfi,
		"ls_option_session_operation_max_search_depth",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetRecordDestination,
		libScrapliFfi,
		"ls_option_session_record_destination",
	)

	// auth
	purego.RegisterLibFunc(
		&m.Options.Auth.SetUsername,
		libScrapliFfi,
		"ls_option_auth_username",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetPassword,
		libScrapliFfi,
		"ls_option_auth_password",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetPrivateKeyPath,
		libScrapliFfi,
		"ls_option_auth_private_key_path",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetPrivateKeyPassphrase,
		libScrapliFfi,
		"ls_option_auth_private_key_passphrase",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetDriverOptionAuthLookupKeyValue,
		libScrapliFfi,
		"ls_option_auth_set_lookup_key_value",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetInSessionAuthBypass,
		libScrapliFfi,
		"ls_option_auth_in_session_auth_bypass",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetUsernamePattern,
		libScrapliFfi,
		"ls_option_auth_username_pattern",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetPasswordPattern,
		libScrapliFfi,
		"ls_option_auth_password_pattern",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetPassphrasePattern,
		libScrapliFfi,
		"ls_option_auth_private_key_passphrase_pattern",
	)

	// transport bin
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetBin,
		libScrapliFfi,
		"ls_option_transport_bin_bin",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetExtraOpenArgs,
		libScrapliFfi,
		"ls_option_transport_bin_extra_open_args",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetOverrideOpenArgs,
		libScrapliFfi,
		"ls_option_transport_bin_override_open_args",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetSSHConfigPath,
		libScrapliFfi,
		"ls_option_transport_bin_ssh_config_path",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetKnownHostsPath,
		libScrapliFfi,
		"ls_option_transport_bin_known_hosts_path",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetEnableStrictKey,
		libScrapliFfi,
		"ls_option_transport_bin_enable_strict_key",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetTermHeight,
		libScrapliFfi,
		"ls_option_transport_bin_term_height",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetTermWidth,
		libScrapliFfi,
		"ls_option_transport_bin_term_width",
	)

	// transport ssh2
	purego.RegisterLibFunc(
		&m.Options.TransportSSH2.SetKnownHostsPath,
		libScrapliFfi,
		"ls_option_transport_ssh2_known_hosts_path",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportSSH2.SetLibSSH2Trace,
		libScrapliFfi,
		"ls_option_transport_ssh2_libssh2trace",
	)

	// transport test
	purego.RegisterLibFunc(
		&m.Options.TransportTest.SetF,
		libScrapliFfi,
		"ls_option_transport_test_f",
	)

	// netconf
	purego.RegisterLibFunc(
		&m.Options.Netconf.SetErrorTag,
		libScrapliFfi,
		"ls_option_netconf_error_tag",
	)

	purego.RegisterLibFunc(
		&m.Options.Netconf.SetPreferredVersion,
		libScrapliFfi,
		"ls_option_netconf_preferred_version",
	)

	purego.RegisterLibFunc(
		&m.Options.Netconf.SetMessagePollIntervalNS,
		libScrapliFfi,
		"ls_option_netconf_message_poll_interval",
	)
}
