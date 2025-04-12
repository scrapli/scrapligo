package ffi

import "github.com/ebitengine/purego"

func registerShared(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.Shared.Free, libScrapliFfi, "freeDriver")
	purego.RegisterLibFunc(&m.Shared.Open, libScrapliFfi, "openDriver")
	purego.RegisterLibFunc(&m.Shared.Close, libScrapliFfi, "closeDriver")

	purego.RegisterLibFunc(&m.Shared.Read, libScrapliFfi, "readSession")
	purego.RegisterLibFunc(&m.Shared.Write, libScrapliFfi, "writeSession")
}

func registerCli(m *Mapping, libScrapliFfi uintptr) {
	// ENHANCEMENT?: is it possible to have my own register funcs that bypass reflection?
	//  driver creation/destruction
	purego.RegisterLibFunc(&m.Cli.Alloc, libScrapliFfi, "allocCliDriver")

	// TODO before we go much further should rationalize the naming of the extern funcs
	//   so probably 1) figure out what is idiomatic for the c abi anyway, and
	//   2) make sure things are "cli" and "netconf" and maybe "shared/common"?
	purego.RegisterLibFunc(&m.Cli.PollOperation, libScrapliFfi, "pollOperation")
	purego.RegisterLibFunc(&m.Cli.FetchOperation, libScrapliFfi, "fetchOperation")

	purego.RegisterLibFunc(&m.Cli.EnterMode, libScrapliFfi, "enterMode")
	purego.RegisterLibFunc(&m.Cli.GetPrompt, libScrapliFfi, "getPrompt")
	purego.RegisterLibFunc(&m.Cli.SendInput, libScrapliFfi, "sendInput")
	purego.RegisterLibFunc(&m.Cli.SendPromptedInput, libScrapliFfi, "sendPromptedInput")
}

func registerNetconf(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.Netconf.Alloc, libScrapliFfi, "allocNetconfDriver")

	purego.RegisterLibFunc(
		&m.Netconf.PollOperation,
		libScrapliFfi,
		"netconfPollOperation",
	)
	purego.RegisterLibFunc(
		&m.Netconf.FetchOperation,
		libScrapliFfi,
		"netconfFetchOperation",
	)

	purego.RegisterLibFunc(&m.Netconf.GetSessionID, libScrapliFfi, "netconfGetSessionID")

	purego.RegisterLibFunc(&m.Netconf.RawRPC, libScrapliFfi, "netconfRawRpc")

	purego.RegisterLibFunc(&m.Netconf.GetConfig, libScrapliFfi, "netconfGetConfig")
	purego.RegisterLibFunc(&m.Netconf.EditConfig, libScrapliFfi, "netconfEditConfig")
	purego.RegisterLibFunc(&m.Netconf.CopyConfig, libScrapliFfi, "netconfCopyConfig")
	purego.RegisterLibFunc(&m.Netconf.DeleteConfig, libScrapliFfi, "netconfDeleteConfig")
	purego.RegisterLibFunc(&m.Netconf.Lock, libScrapliFfi, "netconfLock")
	purego.RegisterLibFunc(&m.Netconf.Unlock, libScrapliFfi, "netconfUnlock")
	purego.RegisterLibFunc(&m.Netconf.Get, libScrapliFfi, "netconfGet")
	purego.RegisterLibFunc(&m.Netconf.CloseSession, libScrapliFfi, "netconfCloseSession")
	purego.RegisterLibFunc(&m.Netconf.KillSession, libScrapliFfi, "netconfKillSession")

	purego.RegisterLibFunc(&m.Netconf.Commit, libScrapliFfi, "netconfCommit")
	purego.RegisterLibFunc(&m.Netconf.Discard, libScrapliFfi, "netconfDiscard")
	purego.RegisterLibFunc(&m.Netconf.CancelCommit, libScrapliFfi, "netconfCancelCommit")
	purego.RegisterLibFunc(&m.Netconf.Validate, libScrapliFfi, "netconfValidate")

	purego.RegisterLibFunc(&m.Netconf.GetSchema, libScrapliFfi, "netconfGetSchema")
	purego.RegisterLibFunc(&m.Netconf.GetData, libScrapliFfi, "netconfGetData")
	purego.RegisterLibFunc(&m.Netconf.EditData, libScrapliFfi, "netconfEditData")
	purego.RegisterLibFunc(&m.Netconf.Action, libScrapliFfi, "netconfAction")
}

func registerOptions(m *Mapping, libScrapliFfi uintptr) {
	// session
	purego.RegisterLibFunc(
		&m.Options.Session.SetReadSize,
		libScrapliFfi,
		"setDriverOptionSessionReadSize",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetReadDelayMinNs,
		libScrapliFfi,
		"setDriverOptionSessionReadDelayMinNs",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetReadDelayMaxNs,
		libScrapliFfi,
		"setDriverOptionSessionReadDelayMaxNs",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetReadDelayBackoffFactor,
		libScrapliFfi,
		"setDriverOptionSessionReadDelayBackoffFactor",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetReturnChar,
		libScrapliFfi,
		"setDriverOptionSessionReturnChar",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetOperationTimeoutNs,
		libScrapliFfi,
		"setDriverOptionSessionOperationTimeoutNs",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetOperationMaxSearchDepth,
		libScrapliFfi,
		"setDriverOptionSessionOperationMaxSearchDepth",
	)
	purego.RegisterLibFunc(
		&m.Options.Session.SetRecorderPath,
		libScrapliFfi,
		"setDriverOptionSessionRecorderPath",
	)

	// auth
	purego.RegisterLibFunc(
		&m.Options.Auth.SetUsername,
		libScrapliFfi,
		"setDriverOptionAuthUsername",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetPassword,
		libScrapliFfi,
		"setDriverOptionAuthPassword",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetPrivateKeyPath,
		libScrapliFfi,
		"setDriverOptionAuthPrivateKeyPath",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetPrivateKeyPassphrase,
		libScrapliFfi,
		"setDriverOptionAuthPrivateKeyPassphrase",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetDriverOptionAuthLookupKeyValue,
		libScrapliFfi,
		"setDriverOptionAuthLookupKeyValue",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetInSessionAuthBypass,
		libScrapliFfi,
		"setDriverOptionAuthInSessionAuthBypass",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetUsernamePattern,
		libScrapliFfi,
		"setDriverOptionAuthUsernamePattern",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetPasswordPattern,
		libScrapliFfi,
		"setDriverOptionAuthPasswordPattern",
	)
	purego.RegisterLibFunc(
		&m.Options.Auth.SetPassphrasePattern,
		libScrapliFfi,
		"setDriverOptionAuthPassphrasePattern",
	)

	// transport bin
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetBin,
		libScrapliFfi,
		"setDriverOptionBinTransportBin",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetExtraOpenArgs,
		libScrapliFfi,
		"setDriverOptionBinTransportExtraOpenArgs",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetOverrideOpenArgs,
		libScrapliFfi,
		"setDriverOptionBinTransportOverrideOpenArgs",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetSSHConfigPath,
		libScrapliFfi,
		"setDriverOptionBinTransportSSHConfigPath",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetKnownHostsPath,
		libScrapliFfi,
		"setDriverOptionBinTransportKnownHostsPath",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetEnableStrictKey,
		libScrapliFfi,
		"setDriverOptionBinTransportEnableStrictKey",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetTermHeight,
		libScrapliFfi,
		"setDriverOptionBinTransportTermHeight",
	)
	purego.RegisterLibFunc(
		&m.Options.TransportBin.SetTermWidth,
		libScrapliFfi,
		"setDriverOptionBinTransportTermWidth",
	)

	// transport ssh2
	purego.RegisterLibFunc(
		&m.Options.TransportSSH2.SetLibSSH2Trace,
		libScrapliFfi,
		"setDriverOptionSSH2TransportSSH2Trace",
	)

	// transport test
	purego.RegisterLibFunc(
		&m.Options.TransportTest.SetF,
		libScrapliFfi,
		// TODO shouldnt this be file? or did i name it etst... i forget
		"setDriverOptionTestTransportF",
	)
}
