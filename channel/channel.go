package channel

import (
	"bytes"
	"errors"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/scrapli/scrapligo/logging"
	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util"
)

const (
	// DefaultTimeoutOpsSeconds is the default time value for operations -- 60 seconds.
	DefaultTimeoutOpsSeconds = 60
	// DefaultReadDelayMilliSeconds is the default value for the delay between reads of the
	// transport -- 5 milliseconds.
	DefaultReadDelayMilliSeconds = 5
	// DefaultReturnChar is the character used to send an "enter" key to the device, "\n".
	DefaultReturnChar = "\n"

	redacted = "redacted"
)

var (
	promptPattern     *regexp.Regexp //nolint:gochecknoglobals
	promptPatternOnce sync.Once      //nolint:gochecknoglobals
)

func getPromptPattern() *regexp.Regexp {
	promptPatternOnce.Do(func() {
		promptPattern = regexp.MustCompile(`(?im)^[a-z\d.\-@()/:]{1,48}[#>$]\s*$`)
	})

	return promptPattern
}

// NewChannel returns a scrapligo Channel object.
func NewChannel(
	l *logging.Instance,
	t *transport.Transport,
	options ...util.Option,
) (*Channel, error) {
	patterns := getAuthPatterns()

	c := &Channel{
		l: l,
		t: t,

		TimeoutOps: DefaultTimeoutOpsSeconds * time.Second,
		ReadDelay:  DefaultReadDelayMilliSeconds * time.Millisecond,

		UsernamePattern:   patterns.username,
		PasswordPattern:   patterns.password,
		PassphrasePattern: patterns.passphrase,

		PromptPattern: getPromptPattern(),
		ReturnChar:    []byte(DefaultReturnChar),

		done: make(chan bool),

		Q:    util.NewQueue(),
		Errs: make(chan error),

		ChannelLog: nil,
	}

	for _, option := range options {
		err := option(c)
		if err != nil {
			if !errors.Is(err, util.ErrIgnoredOption) {
				return nil, err
			}
		}
	}

	return c, nil
}

// Channel is an object that sits "on top" of a scrapligo Transport object, its purpose in life is
// to read data from the transport into its Q, and provide methods to read "until" an input or an
// expected prompt is seen.
type Channel struct {
	l *logging.Instance
	t *transport.Transport

	TimeoutOps time.Duration
	ReadDelay  time.Duration

	AuthBypass bool

	UsernamePattern   *regexp.Regexp
	PasswordPattern   *regexp.Regexp
	PassphrasePattern *regexp.Regexp

	PromptPattern *regexp.Regexp
	ReturnChar    []byte

	done chan bool

	Q    *util.Queue
	Errs chan error

	ChannelLog io.Writer
}

// Open opens the underlying Transport and begins the `read` goroutine, this also kicks off any
// in channel authentication (if necessary).
func (c *Channel) Open() error {
	err := c.t.Open()
	if err != nil {
		c.l.Criticalf("error opening channel, error: %s", err)

		return err
	}

	c.l.Debug("starting channel read loop")

	go c.read()

	if !c.AuthBypass {
		var b []byte

		switch tt := c.t.Impl.(type) {
		case *transport.System:
			c.l.Debug("transport type is 'system', begin in channel ssh authentication")

			b, err = c.AuthenticateSSH([]byte(c.t.Args.Password), []byte(tt.SSHArgs.PrivateKeyPassPhrase))
			if err != nil {
				return err
			}
		case *transport.Telnet:
			c.l.Debug("transport type is 'telnet', begin in channel telnet authentication")

			b, err = c.AuthenticateTelnet([]byte(c.t.Args.User), []byte(c.t.Args.Password))
			if err != nil {
				return err
			}
		}

		if len(b) > 0 {
			// requeue any buffer data we get during in channel authentication back onto the
			// read buffer. mostly this should only be relevant for netconf where we need to
			// read the server capabilities.
			c.Q.Requeue(b)
		}
	} else {
		c.l.Debug("auth bypass is enabled, skipping in channel auth check")
	}

	return nil
}

// Close signals to stop the channel read loop and closes the underlying Transport object.
func (c *Channel) Close() error {
	c.done <- true

	return c.t.Close()
}

type result struct {
	b   []byte
	err error
}

func (c *Channel) processOut(b []byte, strip bool) []byte {
	lines := bytes.Split(b, []byte("\n"))

	cleanLines := make([][]byte, len(lines))
	for i, l := range lines {
		cleanLines[i] = bytes.TrimRight(l, " ")
	}

	b = bytes.Join(cleanLines, []byte("\n"))

	if strip {
		b = c.PromptPattern.ReplaceAll(b, nil)
	}

	b = bytes.Trim(b, string(c.ReturnChar))
	b = bytes.Trim(b, "\n")

	return b
}

// GetTimeout returns the target timeout for an operation based on the TimeoutOps attribute of the
// Channel and the value t.
func (c *Channel) GetTimeout(t time.Duration) time.Duration {
	if t == -1 {
		return c.TimeoutOps
	}

	if t == 0 {
		return util.MaxTimeout * time.Second
	}

	return t
}
