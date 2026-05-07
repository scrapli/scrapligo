package ffi

import "github.com/ebitengine/purego"

// Mapping holds mappings to the libscrapli external functions.
type Mapping struct {
	// AssertNoLeaks returns "true" if no leaks were found, otherwise false.
	AssertNoLeaks func() bool

	Shared  SharedMapping
	Session SessionMapping
	Cli     CliMapping
	Netconf NetconfMapping
}

func registerShared(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.Shared.GetPollFd, libScrapliFfi, "ls_shared_get_poll_fd")
	purego.RegisterLibFunc(&m.Shared.Free, libScrapliFfi, "ls_shared_free")

	purego.RegisterLibFunc(&m.Shared.AllocDriverOptions, libScrapliFfi, "ls_alloc_driver_options")
	purego.RegisterLibFunc(&m.Shared.FreeDriverOptions, libScrapliFfi, "ls_free_driver_options")

	purego.RegisterLibFunc(&m.Shared.FetchOptionsSize, libScrapliFfi, "ls_fetch_options_size")
	purego.RegisterLibFunc(&m.Shared.FetchOptions, libScrapliFfi, "ls_fetch_options")
}

// SharedMapping holds common mappings for both cli and netconf drivers.
type SharedMapping struct {
	GetPollFd func(driverPtr uintptr) uint32

	// Free releases the memory of the driver object at driverPtr -- this should be called *after*
	// Close where possible.
	Free func(driverPtr uintptr)

	AllocDriverOptions func() uintptr
	FreeDriverOptions  func(optionsPtr uintptr)

	FetchOptionsSize func(
		optionsPtr uintptr,
		optionsSize *uintptr,
	) uintptr
	FetchOptions func(
		optionsPtr uintptr,
		options *[]byte,
	)
}
