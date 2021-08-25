package base

import (
	"errors"
	"regexp"
	"time"

	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/transport"
)

func NewChannel(
	t *transport.Transport,
	options ...Option,
) (*channel.Channel, error) {
	c := &channel.Channel{
		Transport:              t,
		CommsPromptPattern:     regexp.MustCompile(`(?im)^[a-z0-9.\-@()/:]{1,48}[#>$]\s*$`),
		CommsReturnChar:        "\n",
		CommsPromptSearchDepth: 1000,
		TimeoutOps:             60 * time.Second,
		Host:                   t.BaseTransportArgs.Host,
		Port:                   t.BaseTransportArgs.Port,
	}

	for _, option := range options {
		err := option(c)

		if err != nil {
			if errors.Is(err, ErrIgnoredOption) {
				continue
			} else {
				return nil, err
			}
		}
	}

	return c, nil
}
