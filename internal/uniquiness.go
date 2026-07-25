package internal

import (
	"sync"
)

var (
	userDataDispatcherInst     *userDataDispatcher //nolint: gochecknoglobals
	userDataDispatcherInstOnce sync.Once           //nolint: gochecknoglobals
)

// UserDataDispatcher is the interface returned when fetching the singleton userDataDispatcher. The
// dispatcher exists to generate a unique id per connection -- this id is then used in the logger
// singleton (and maybe other places? netconf caps? read w/ callbacks?) to associate a Cli/Netconf
// object to a slot on the singleton.
type UserDataDispatcher interface {
	Register() uintptr
}

// GetUserDataDispatcherr returns the UserDataDispatcher singleton.
func GetUserDataDispatcherr() UserDataDispatcher { //nolint: ireturn
	userDataDispatcherInstOnce.Do(
		func() {
			userDataDispatcherInst = &userDataDispatcher{
				lock: sync.Mutex{},
			}
		},
	)

	return userDataDispatcherInst
}

type userDataDispatcher struct {
	lock        sync.Mutex
	curUserData uintptr
}

func (u *userDataDispatcher) Register() uintptr {
	u.lock.Lock()
	defer u.lock.Unlock()

	userData := u.curUserData
	u.curUserData++

	return userData
}
