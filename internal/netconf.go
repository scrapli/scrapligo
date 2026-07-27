package internal

import (
	"strings"
	"sync"

	"github.com/ebitengine/purego"
)

var (
	netconfCapabiltiesDispatcherInst     *netconfCapabilitiesDispatcher //nolint: gochecknoglobals
	netconfCapabiltiesDispatcherInstOnce sync.Once                      //nolint: gochecknoglobals
)

type capabilitiesCallbackF = func(caps string) string

// NetconfCapabiltiesDispatcher is the interface returned when fetching the singleton
// netconfCapabilitiesDispatcher. Same story sas logger and session recorder dispatcher.
type NetconfCapabiltiesDispatcher interface {
	Register(userData uintptr, f capabilitiesCallbackF)
	Deregister(userData uintptr)

	GetCapabilitiesCallback() uintptr
}

// GetNetconfCapabiltiesDispatcher returns the NetconfCapabiltiesDispatcher singleton.
func GetNetconfCapabiltiesDispatcher() NetconfCapabiltiesDispatcher { //nolint: ireturn
	netconfCapabiltiesDispatcherInstOnce.Do(
		func() {
			netconfCapabiltiesDispatcherInst = &netconfCapabilitiesDispatcher{
				lock:                  sync.RWMutex{},
				capabilitiesCallbacks: map[uintptr]capabilitiesCallbackF{},
			}

			cb := purego.NewCallback(netconfCapabiltiesDispatcherInst.capabilities)

			netconfCapabiltiesDispatcherInst.cb = cb
		},
	)

	return netconfCapabiltiesDispatcherInst
}

type netconfCapabilitiesDispatcher struct {
	lock                  sync.RWMutex
	capabilitiesCallbacks map[uintptr]capabilitiesCallbackF
	cb                    uintptr
}

func (n *netconfCapabilitiesDispatcher) Register(userData uintptr, f capabilitiesCallbackF) {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.capabilitiesCallbacks[userData] = f
}

func (n *netconfCapabilitiesDispatcher) Deregister(userData uintptr) {
	n.lock.Lock()
	defer n.lock.Unlock()

	delete(n.capabilitiesCallbacks, userData)
}

func (n *netconfCapabilitiesDispatcher) GetCapabilitiesCallback() uintptr {
	return n.cb
}

func (n *netconfCapabilitiesDispatcher) capabilities(userData uintptr, message *string) *string {
	n.lock.RLock()
	defer n.lock.RUnlock()

	var out string

	cf, ok := n.capabilitiesCallbacks[userData]
	if !ok {
		return &out
	}

	msg := strings.Clone(*message)

	out = cf(msg)

	return &out
}
