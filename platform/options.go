package platform

import (
	"regexp"
	"time"

	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/util"
)

const (
	port = "port"

	authBypass    = "auth-bypass"
	authStrictKey = "auth-strict-key"

	promptPattern     = "prompt-pattern"
	usernamePattern   = "username-pattern"
	passwordPattern   = "password-pattern"
	passphrasePattern = "passphrase-pattern"

	returnChar = "return-char"

	// read delay in seconds for channel read loop.
	readDelay = "read-delay"

	// timeouts in seconds.
	timeoutOps       = "timeout-ops"
	timeoutTransport = "timeout-transport"

	transportType = "transport-type"
	// read size for transport read chunk.
	transportReadSize  = "read-size"
	transportPtyHeight = "transport-pty-height"
	transportPtyWidth  = "transport-pty-width"

	transportSystemOpenArgs = "transport-system-open-args"
)

type optionDefinition struct {
	Option string      `json:"option" yaml:"option"`
	Value  interface{} `json:"value"  yaml:"value"`
}

type optionDefinitions []*optionDefinition

func (o *optionDefinitions) asOptions() []util.Option { //nolint: gocyclo
	opts := make([]util.Option, len(*o))

	for i, opt := range *o {
		switch opt.Option {
		case port:
			opts[i] = options.WithPort(opt.Value.(int))
		case authBypass:
			opts[i] = options.WithAuthBypass()
		case authStrictKey:
			opts[i] = options.WithAuthNoStrictKey()
		case promptPattern:
			opts[i] = options.WithPromptPattern(regexp.MustCompile(opt.Value.(string)))
		case usernamePattern:
			opts[i] = options.WithUsernamePattern(regexp.MustCompile(opt.Value.(string)))
		case passwordPattern:
			opts[i] = options.WithPasswordPattern(regexp.MustCompile(opt.Value.(string)))
		case passphrasePattern:
			opts[i] = options.WithPassphrasePattern(regexp.MustCompile(opt.Value.(string)))
		case returnChar:
			opts[i] = options.WithReturnChar(opt.Value.(string))
		case readDelay:
			opts[i] = options.WithReadDelay(
				time.Duration(opt.Value.(float64) * float64(time.Second)),
			)
		case timeoutOps:
			opts[i] = options.WithTimeoutOps(
				time.Duration(opt.Value.(float64) * float64(time.Second)),
			)
		case timeoutTransport:
			opts[i] = options.WithTimeoutTransport(
				time.Duration(opt.Value.(float64) * float64(time.Second)),
			)
		case transportType:
			opts[i] = options.WithTransportType(opt.Value.(string))
		case transportReadSize:
			opts[i] = options.WithTransportReadSize(opt.Value.(int))
		case transportPtyHeight:
			opts[i] = options.WithTermHeight(opt.Value.(int))
		case transportPtyWidth:
			opts[i] = options.WithTermWidth(opt.Value.(int))
		case transportSystemOpenArgs:
			opts[i] = options.WithSystemTransportOpenArgs(opt.Value.([]string))
		}
	}

	return opts
}
