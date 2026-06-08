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

	purego.RegisterLibFunc(&m.Shared.fetchOptionsSize, libScrapliFfi, "ls_fetch_options_size")
	purego.RegisterLibFunc(&m.Shared.fetchOptions, libScrapliFfi, "ls_fetch_options")
}

// SharedMapping holds common mappings for both cli and netconf drivers.
type SharedMapping struct {
	GetPollFd func(driverPtr uintptr) uint32

	// Free releases the memory of the driver object at driverPtr -- this should be called *after*
	// Close where possible.
	Free func(driverPtr uintptr)

	AllocDriverOptions func() uintptr
	FreeDriverOptions  func(optionsPtr uintptr)

	fetchOptionsSize func(
		optionsPtr uintptr,
		optionsSize *uintptr,
	) uint8
	fetchOptions func(
		optionsPtr uintptr,
		options *[]byte,
	) uint8
}

// FetchOptionsSize fetches the size of the options json representation for the options at
// optionsPtr.
func (m *SharedMapping) FetchOptionsSize(
	optionsPtr uintptr,
	optionsSize *uintptr,
) error {
	return newLibScrapliResult(
		m.fetchOptionsSize(
			optionsPtr,
			optionsSize,
		),
		"fetch options size failed",
	).check()
}

// FetchOptions fetches the options as json string, writing it into the options slice ptr which
// must be pre allocated (hint: FetchOptionsSize).
func (m *SharedMapping) FetchOptions(
	optionsPtr uintptr,
	options *[]byte,
) error {
	return newLibScrapliResult(
		m.fetchOptions(
			optionsPtr,
			options,
		),
		"fetch options failed",
	).check()
}
