package ffi

import "github.com/ebitengine/purego"

func registerDriver(m *Mapping, libScrapliFfi uintptr) {
	// TODO is it possible to have my own register funcs that bypass reflection?
	//  driver creation/destruction
	purego.RegisterLibFunc(&m.Driver.AllocFromYaml, libScrapliFfi, "allocDriverFromYaml")
	purego.RegisterLibFunc(
		&m.Driver.AllocFromYamlString,
		libScrapliFfi,
		"allocDriverFromYamlString",
	)
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
