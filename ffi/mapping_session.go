package ffi

import "github.com/ebitengine/purego"

func registerSession(m *Mapping, libScrapliFfi uintptr) {
	purego.RegisterLibFunc(&m.Session.read, libScrapliFfi, "ls_session_read")
	purego.RegisterLibFunc(&m.Session.write, libScrapliFfi, "ls_session_write")
	purego.RegisterLibFunc(&m.Session.writeAndReturn, libScrapliFfi, "ls_session_write_and_return")
	purego.RegisterLibFunc(&m.Session.writeReturn, libScrapliFfi, "ls_session_write_return")
}

// SessionMapping holds session specific mappings.
type SessionMapping struct {
	read func(
		driverPtr uintptr,
		buf *[]byte,
		readSize *uintptr,
	) uint8
	write func(
		driverPtr uintptr,
		buf string,
		redacted bool,
	) uint8
	writeAndReturn func(
		driverPtr uintptr,
		buf string,
		redacted bool,
	) uint8
	writeReturn func(
		driverPtr uintptr,
	) uint8
}

// Read exposes the libscrapli session read function.
func (m *SessionMapping) Read(
	driverPtr uintptr,
	buf *[]byte,
	readSize *uintptr,
) error {
	return newLibScrapliResult(
		m.read(driverPtr, buf, readSize),
		"failed executing read",
	).check()
}

// Write exposes the libscrapli session write function.
func (m *SessionMapping) Write(
	driverPtr uintptr,
	buf string,
	redacted bool,
) error {
	return newLibScrapliResult(
		m.write(driverPtr, buf, redacted),
		"failed executing write",
	).check()
}

// WriteAndReturn exposes the libscrapli session writeAndReturn function.
func (m *SessionMapping) WriteAndReturn(
	driverPtr uintptr,
	buf string,
	redacted bool,
) error {
	return newLibScrapliResult(
		m.writeAndReturn(driverPtr, buf, redacted),
		"failed executing writeAndReturn",
	).check()
}

// WriteReturn exposes the libscrapli session writeReturn function.
func (m *SessionMapping) WriteReturn(
	driverPtr uintptr,
) error {
	return newLibScrapliResult(
		m.writeReturn(driverPtr),
		"failed executing writeReturn",
	).check()
}
