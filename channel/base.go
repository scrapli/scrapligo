package channel

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/transport"
)

// ErrAuthTimeout error for channel auth timeouts.
var ErrAuthTimeout = errors.New("channel authentication timed out")

// ErrAuthFailedPassword raised when password prompt seen too many times during in channel auth.
var ErrAuthFailedPassword = errors.New(
	"password prompt seen more than twice, assuming auth failed",
)

// ErrAuthFailedPassphrase raised when passphrase prompt seen too many times during in channel auth.
var ErrAuthFailedPassphrase = errors.New(
	"passphrase prompt seen more than twice, assuming auth failed",
)

// ErrChannelTimeout error for channel operation timeouts.
var ErrChannelTimeout = errors.New("channel operation timed out")

const (
	loginSeenMax      = 2
	passwordSeenMax   = 2
	passphraseSeenMax = 2
	// MaxTimeout maximum allowable timeout value -- one day.
	MaxTimeout = 86_400
)

// Channel struct representing the channel object.
type Channel struct {
	Transport              *transport.Transport
	CommsPromptPattern     *regexp.Regexp
	CommsReturnChar        string
	CommsPromptSearchDepth int
	TimeoutOps             time.Duration
	Host                   string
	Port                   int
	ChannelLog             io.Writer
}

type channelResult struct {
	result []byte
	error  error
}

// Write writes bytes input into the channel, redacted (currently unused) signals that the input
// should not be written in the log output.
func (c *Channel) Write(channelInput []byte, redacted bool) error {
	logOutput := string(channelInput)
	if redacted {
		logOutput = "REDACTED"
	}

	logging.LogDebug(c.FormatLogMessage("write", fmt.Sprintf("write: %s", logOutput)))

	err := c.Transport.Write(channelInput)

	return err
}

// SendReturn convenience function to send the return character.
func (c *Channel) SendReturn() error {
	return c.Write([]byte(c.CommsReturnChar), false)
}

// WriteAndReturn convenience function to write input and send the return character.
func (c *Channel) WriteAndReturn(channelInput []byte, redacted bool) error {
	err := c.Write(channelInput, redacted)
	if err != nil {
		return err
	}

	err = c.SendReturn()
	if err != nil {
		return err
	}

	return nil
}

func (c *Channel) readUntilInput(channelInput []byte) ([]byte, error) {
	var b []byte

	if len(channelInput) == 0 {
		return b, nil
	}

	for {
		chunk, err := c.Read()
		b = append(b, chunk...)

		if err != nil {
			return b, err
		}

		if bytes.Contains(b, channelInput) {
			return b, err
		}
	}
}

func (c *Channel) readUntilPrompt() ([]byte, error) {
	matchPattern := c.CommsPromptPattern

	var b []byte

	for {
		chunk, err := c.Read()
		b = append(b, chunk...)

		if err != nil {
			return b, err
		}

		channelMatch := matchPattern.Match(b)
		if channelMatch {
			logging.LogDebug(c.FormatLogMessage("debug", "found prompt match"))

			return b, err
		}
	}
}

func (c *Channel) readUntilExplicitPrompt(prompts []*regexp.Regexp) ([]byte, error) {
	var b []byte

	for {
		chunk, err := c.Read()
		b = append(b, chunk...)

		if err != nil {
			return b, err
		}

		for _, prompt := range prompts {
			channelMatch := prompt.Match(b)

			if channelMatch {
				logging.LogDebug(c.FormatLogMessage("debug", "found prompt match"))

				return b, err
			}
		}
	}
}

// Read reads bytes off the transport, handles some basic "massaging" of data to remove null bytes,
// \r characters, as well as stripping out any ANSI characters in the output.
func (c *Channel) Read() ([]byte, error) {
	chunk, err := c.Transport.Read()

	// is there ever a time when we should *not* replace *all* null bytes? previously this was just
	// a trim, but for some connections (nxos telnet in particular) null byte would sneak in and
	// not be leading or trailing... it would cause chaos....
	b := bytes.ReplaceAll(chunk, []byte("\x00"), []byte(""))
	b = bytes.ReplaceAll(b, []byte("\r"), []byte(""))

	if bytes.Contains(b, []byte("\x1b")) {
		logging.LogDebug(c.FormatLogMessage("debug", "stripping ansi chars..."))

		b = util.StripAnsi(b)
	}

	logging.LogDebug(c.FormatLogMessage("debug", fmt.Sprintf("read: %s", b)))

	if c.ChannelLog != nil {
		_, channelLogErr := c.ChannelLog.Write(b)

		if channelLogErr != nil {
			logging.LogError(c.FormatLogMessage("error", "error writing to channel log"))
		}
	}

	return b, err
}

// RestructureOutput strips prompt (if necessary) from output and trim any null space.
func (c *Channel) RestructureOutput(output []byte, stripPrompt bool) []byte {
	outputLines := bytes.Split(output, []byte("\n"))

	cleanLines := make([][]byte, 0)

	for _, l := range outputLines {
		cl := bytes.TrimRight(l, " ")
		cleanLines = append(cleanLines, cl)
	}

	cleanOutput := bytes.Join(cleanLines, []byte("\n"))

	if stripPrompt {
		cleanOutput = c.CommsPromptPattern.ReplaceAll(cleanOutput, []byte(""))
	}

	cleanOutput = bytes.TrimLeft(cleanOutput, c.CommsReturnChar)
	cleanOutput = bytes.TrimRight(cleanOutput, " ")
	cleanOutput = bytes.TrimRight(cleanOutput, "\n")

	return cleanOutput
}

// DetermineOperationTimeout determines timeout to use for channel operation.
func (c *Channel) DetermineOperationTimeout(timeoutOps time.Duration) time.Duration {
	opTimeout := c.TimeoutOps

	if timeoutOps > 0 {
		opTimeout = timeoutOps
	}

	if opTimeout <= 0 {
		opTimeout = MaxTimeout * time.Second
	}

	return opTimeout
}

// FormatLogMessage formats log message payload, adding contextual info about the host.
func (c *Channel) FormatLogMessage(level, msg string) string {
	return logging.FormatLogMessage(level, c.Host, c.Port, msg)
}
