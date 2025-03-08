package ffi

import "github.com/ebitengine/purego"

func registerDriver(m *Mapping, libScrapliFfi uintptr) {
	// TODO is it possible to have my own register funcs that bypass reflection?
	//  driver creation/destruction
	purego.RegisterLibFunc(&m.Driver.Alloc, libScrapliFfi, "allocDriver")
	purego.RegisterLibFunc(&m.Driver.Free, libScrapliFfi, "freeDriver")

	purego.RegisterLibFunc(&m.Driver.Open, libScrapliFfi, "openDriver")
	purego.RegisterLibFunc(&m.Driver.Close, libScrapliFfi, "closeDriver")

	purego.RegisterLibFunc(&m.Driver.PollOperation, libScrapliFfi, "pollOperation")
	purego.RegisterLibFunc(&m.Driver.FetchOperation, libScrapliFfi, "fetchOperation")

	purego.RegisterLibFunc(&m.Driver.EnterMode, libScrapliFfi, "enterMode")
	purego.RegisterLibFunc(&m.Driver.GetPrompt, libScrapliFfi, "getPrompt")
	purego.RegisterLibFunc(&m.Driver.SendInput, libScrapliFfi, "sendInput")
	purego.RegisterLibFunc(&m.Driver.SendPromptedInput, libScrapliFfi, "sendPromptedInput")
}

func registerNetconf(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.DriverNetconf.Alloc, libScrapliFfi, "netconfAllocDriver")
	purego.RegisterLibFunc(&m.DriverNetconf.Free, libScrapliFfi, "netconfFreeDriver")

	purego.RegisterLibFunc(&m.DriverNetconf.Open, libScrapliFfi, "netconfOpenDriver")
	purego.RegisterLibFunc(&m.DriverNetconf.Close, libScrapliFfi, "netconfCloseDriver")

	purego.RegisterLibFunc(
		&m.DriverNetconf.PollOperation,
		libScrapliFfi,
		"netconfPollOperation",
	)
	purego.RegisterLibFunc(
		&m.DriverNetconf.FetchOperation,
		libScrapliFfi,
		"netconfFetchOperation",
	)

	purego.RegisterLibFunc(&m.DriverNetconf.GetConfig, libScrapliFfi, "netconfGetConfig")
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
}
