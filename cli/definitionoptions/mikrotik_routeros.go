package definitionoptions

import (
	"fmt"

	scrapligointernal "github.com/scrapli/scrapligo/internal"
	scrapligooptions "github.com/scrapli/scrapligo/options"
)

const mikrotikRouterOS = "mikrotik_routeros"

// WithMikrotikRouterOSUsername is a static option that appends "+tc" to a users configured username
// -- this disables terminal auto detection and terminal colors.
func WithMikrotikRouterOSUsername() scrapligooptions.Option {
	return func(o *scrapligointernal.Options) error {
		o.Auth.Username = fmt.Sprintf("%s+tc", o.Auth.Username)

		return nil
	}
}

// WithMikrotikRouterOSReturnChar sets the return char to \r\n.
func WithMikrotikRouterOSReturnChar() scrapligooptions.Option {
	return func(o *scrapligointernal.Options) error {
		o.Session.ReturnChar = "\r\n"

		return nil
	}
}

func registerMikrotikRouterOSOptions() []scrapligooptions.Option {
	return []scrapligooptions.Option{
		WithMikrotikRouterOSUsername(),
		WithMikrotikRouterOSReturnChar(),
	}
}
