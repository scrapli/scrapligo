// Package definitionoptions
// DO NOT EDIT, GENERATED FILE
package definitionoptions

import (
	"sync"

	scrapligooptions "github.com/scrapli/scrapligo/v2/options"
)

var (
	platformOptionsInst      *platformOptions //nolint: gochecknoglobals
	platformOptionsInsttOnce sync.Once        //nolint: gochecknoglobals
)

// PlatformOptions is an interface defining the platform registration/mapping singleton.
type PlatformOptions interface {
	RegisterOptionsForPlatform(s string, opts ...scrapligooptions.Option)
	ClearOptionsForPlatform(s string)
	GetOptionsForPlatform(s string) []scrapligooptions.Option
}

// GetPlatformOptions returns the PlatformOptions singleton.
func GetPlatformOptions() PlatformOptions { //nolint: ireturn
	platformOptionsInsttOnce.Do(func() {
		platformOptionsInst = &platformOptions{
			optLock: &sync.Mutex{},
			optMap: map[string][]scrapligooptions.Option{
				mikrotikRouterOS: registerMikrotikRouterOSOptions(),
			},
		}
	})

	return platformOptionsInst
}

type platformOptions struct {
	optLock *sync.Mutex
	optMap  map[string][]scrapligooptions.Option
}

// RegisterOptionsForPlatform registers the list of "static" options for the given platform.
func (o *platformOptions) RegisterOptionsForPlatform(
	s string,
	opts ...scrapligooptions.Option,
) {
	o.optLock.Lock()
	defer o.optLock.Unlock()

	o.optMap[s] = append(o.optMap[s], opts...)
}

// ClearOptionsForPlatform clears any loaded options for the given platform.
func (o *platformOptions) ClearOptionsForPlatform(s string) {
	o.optLock.Lock()
	defer o.optLock.Unlock()

	delete(o.optMap, s)
}

// GetOptionsForPlatform returns any "static" options for the given platform.
func (o *platformOptions) GetOptionsForPlatform(s string) []scrapligooptions.Option {
	o.optLock.Lock()
	defer o.optLock.Unlock()

	return o.optMap[s]
}
