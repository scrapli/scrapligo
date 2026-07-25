package internal

import (
	"strings"
	"sync"

	"github.com/ebitengine/purego"
)

var (
	recorderDispatcherInst     *recorderDispatcher //nolint: gochecknoglobals
	recorderDispatcherInstOnce sync.Once           //nolint: gochecknoglobals
)

// RecorderDispatcher is the interface returned when fetching the singleton recorderDispatcher. The
// dispatcher exists because in purego we have a max of 2000 callbacks -- same story as the logging
// dispatcher. While probably less important duplicating the pattern for the recorder is easy
// enough, so here we are.
type RecorderDispatcher interface {
	Register(userData uintptr, f recorderCallbackF)
	Deregister(userData uintptr)

	GetRecorderCallback() uintptr
}

type recorderCallbackF func(s string)

// GetRecorderDispatcher returns the RecorderDispatcher singleton.
func GetRecorderDispatcher() RecorderDispatcher { //nolint: ireturn
	recorderDispatcherInstOnce.Do(
		func() {
			recorderDispatcherInst = &recorderDispatcher{
				lock: sync.RWMutex{},
			}

			cb := purego.NewCallback(recorderDispatcherInst.record)

			recorderDispatcherInst.cb = cb
		},
	)

	return recorderDispatcherInst
}

type recorderDispatcher struct {
	lock      sync.RWMutex
	recorders map[uintptr]recorderCallbackF
	cb        uintptr
}

func (r *recorderDispatcher) Register(userData uintptr, f recorderCallbackF) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.recorders[userData] = f
}

func (r *recorderDispatcher) Deregister(userData uintptr) {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.recorders, userData)
}

func (r *recorderDispatcher) GetRecorderCallback() uintptr {
	return r.cb
}

func (r *recorderDispatcher) record(userData uintptr, message *string) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	rf, ok := r.recorders[userData]
	if !ok {
		return
	}

	msg := strings.Clone(*message)

	rf(msg)
}
