package options

import (
	"bytes"
	"encoding/xml"

	scrapligointernal "github.com/scrapli/scrapligo/v2/internal"
)

// NetconfVersion is an enumish type representing a netconf version (1.0 or 1.1).
type NetconfVersion string

const (
	// NetconfVersion10 represents the netconf version 1.0.
	NetconfVersion10 NetconfVersion = "1.0"
	// NetconfVersion11 represents the netconf version 1.1.
	NetconfVersion11 NetconfVersion = "1.1"
)

// WithNetconfErrorTag sets the error tag substring for a netconf object.
func WithNetconfErrorTag(s string) Option {
	return func(o *scrapligointernal.Options) error {
		o.Netconf.ErrorTag = s

		return nil
	}
}

// WithNetconfPreferredVersion sets the preferred version for a netconf object.
func WithNetconfPreferredVersion(v NetconfVersion) Option {
	return func(o *scrapligointernal.Options) error {
		o.Netconf.PreferredVersion = string(v)

		return nil
	}
}

// WithNetconfMessagePollIntervalNS sets the message poll interval for a netconf object.
func WithNetconfMessagePollIntervalNS(v uint64) Option {
	return func(o *scrapligointernal.Options) error {
		o.Netconf.MessagePollIntervalNS = v

		return nil
	}
}

// WithNetconfCapabilitiesCallback sets the callback function to receive the servers capabilities
// and return custom client capabilities.
func WithNetconfCapabilitiesCallback(f func(serverCapabilities []string) []string) Option {
	return func(o *scrapligointernal.Options) error {
		// this closure *can* panic, it really *should not* ever happen... but because we are in
		// a closure with no way to propagate the error, and any failure here would represent
		// something unrecoverable anyway -- as in the servers capabilities were unmarshal-able,
		// or we somehow couldnt encode the users provided capabilities. this doesn't *feel* great,
		// but feels like its likely not super common use case, and the likelihood is so low that
		// we'll let it happen...
		o.Netconf.CapabilitiesCallback = func(serverHello *string) *string {
			hello := &scrapligointernal.NetconfServerHello{}

			err := xml.Unmarshal([]byte(*serverHello), &hello)
			if err != nil {
				panic(err)
			}

			userCapabilities := f(hello.Capabilities.Capability)

			buf := bytes.NewBufferString("")
			enc := xml.NewEncoder(buf)

			for _, userCapability := range userCapabilities {
				err = enc.EncodeElement(
					userCapability,
					xml.StartElement{
						Name: xml.Name{
							Local: "capability",
						},
					},
				)
				if err != nil {
					panic(err)
				}
			}

			outCaps := buf.String()

			return &outCaps
		}

		return nil
	}
}
