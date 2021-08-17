package channel

import (
	"bytes"
	"regexp"
	"time"

	"github.com/scrapli/scrapligo/logging"
)

type authPatterns struct {
	telnetLoginPattern *regexp.Regexp
	passwordPattern    *regexp.Regexp
	passphrasePattern  *regexp.Regexp
}

var authPatternsInstance *authPatterns //nolint:gochecknoglobals

func getAuthPatterns() *authPatterns {
	if authPatternsInstance == nil {
		authPatternsInstance = &authPatterns{
			telnetLoginPattern: regexp.MustCompile(`(?im)^(.*username:)|(.*login:)\s?$`),
			passwordPattern:    regexp.MustCompile(`(?im)^(.*@.*)?password:\s?$`),
			passphrasePattern:  regexp.MustCompile(`(?i)enter passphrase for key`),
		}
	}

	return authPatternsInstance
}

func (c *Channel) authenticateSSH(authPassword, authPassphrase []byte) *channelResult {
	logging.LogDebug(c.FormatLogMessage("debug", "attempting in channel ssh authentication"))

	patterns := getAuthPatterns()
	passwordPattern := patterns.passwordPattern
	passphrasePattern := patterns.passphrasePattern

	var passwordCount = 0

	var passphraseCount = 0

	var b []byte

	for {
		chunk, err := c.Read()

		if err != nil {
			return &channelResult{
				result: b,
				error:  err,
			}
		}

		processedChunk := bytes.ToLower(bytes.Trim(chunk, "\x00"))

		b = append(b, processedChunk...)

		passwordMatch := passwordPattern.Match(b)
		passphraseMatch := passphrasePattern.Match(b)

		if passwordMatch { //nolint:nestif
			b = []byte{}
			passwordCount++

			if passwordCount > passwordSeenMax {
				return &channelResult{
					result: []byte{},
					error:  ErrAuthFailedPassword,
				}
			}

			logging.LogDebug(c.FormatLogMessage("debug", "found password prompt, sending password"))

			err = c.WriteAndReturn(authPassword, true)
			if err != nil {
				return &channelResult{
					result: []byte{},
					error:  err,
				}
			}
		} else if passphraseMatch {
			b = []byte{}
			passphraseCount++

			if passwordCount > passphraseSeenMax {
				return &channelResult{
					result: []byte{},
					error:  ErrAuthFailedPassphrase,
				}
			}

			logging.LogDebug(c.FormatLogMessage("debug", "found passphrase prompt, sending passphrase"))

			err = c.WriteAndReturn(authPassphrase, true)
			if err != nil {
				return &channelResult{
					result: []byte{},
					error:  err,
				}
			}
		}

		promptMatch := c.CommsPromptPattern.Match(b)
		if promptMatch {
			logging.LogDebug(c.FormatLogMessage("debug", "ssh authentication complete"))

			return &channelResult{
				result: b,
				error:  nil,
			}
		}
	}
}

// AuthenticateSSH Handles "in channel" SSH authentication (for "system" transport).
func (c *Channel) AuthenticateSSH(authPassword, authPassphrase string) ([]byte, error) {
	var _c = make(chan *channelResult)

	go func() {
		r := c.authenticateSSH([]byte(authPassword), []byte(authPassphrase))
		_c <- r
		close(_c)
	}()

	timer := time.NewTimer(c.DetermineOperationTimeout(*c.TimeoutOps))

	select {
	case r := <-_c:
		return r.result, r.error
	case <-timer.C:
		logging.LogError(c.FormatLogMessage("error", "timed out during ssh authentication"))

		return []byte{}, ErrAuthTimeout
	}
}

func (c *Channel) authenticateTelnet(authUsername, authPassword []byte) *channelResult {
	logging.LogDebug(c.FormatLogMessage("debug", "attempting in channel telnet authentication"))

	patterns := getAuthPatterns()
	usernamePattern := patterns.telnetLoginPattern
	passwordPattern := patterns.passwordPattern

	var usernameCount = 0

	var passwordCount = 0

	var b []byte

	for {
		chunk, err := c.Read()

		if err != nil {
			return &channelResult{
				result: b,
				error:  err,
			}
		}

		processedChunk := bytes.ToLower(bytes.Trim(chunk, "\x00"))

		b = append(b, processedChunk...)

		usernameMatch := usernamePattern.Match(b)
		passwordMatch := passwordPattern.Match(b)

		if usernameMatch { //nolint:nestif
			b = []byte{}
			usernameCount++

			if usernameCount > loginSeenMax {
				return &channelResult{
					result: []byte{},
					error:  ErrAuthFailedPassword,
				}
			}

			logging.LogDebug(c.FormatLogMessage("debug", "found login prompt, sending username"))

			err = c.WriteAndReturn(authUsername, false)
			if err != nil {
				return &channelResult{
					result: []byte{},
					error:  err,
				}
			}
		} else if passwordMatch {
			b = []byte{}
			passwordCount++

			if passwordCount > passwordSeenMax {
				return &channelResult{
					result: []byte{},
					error:  ErrAuthFailedPassword,
				}
			}

			logging.LogDebug(c.FormatLogMessage("debug", "found password prompt, sending password"))

			err = c.WriteAndReturn(authPassword, true)
			if err != nil {
				return &channelResult{
					result: []byte{},
					error:  err,
				}
			}
		}

		promptMatch := c.CommsPromptPattern.Match(b)
		if promptMatch {
			logging.LogDebug(c.FormatLogMessage("debug", "telnet authentication complete"))

			return &channelResult{
				result: b,
				error:  nil,
			}
		}
	}
}

// AuthenticateTelnet Handles "in channel" Telnet authentication (for "telnet" transport).
func (c *Channel) AuthenticateTelnet(authUsername, authPassword string) ([]byte, error) {
	var _c = make(chan *channelResult)

	go func() {
		r := c.authenticateTelnet([]byte(authUsername), []byte(authPassword))
		_c <- r
		close(_c)
	}()

	timer := time.NewTimer(c.DetermineOperationTimeout(*c.TimeoutOps))

	select {
	case r := <-_c:
		return r.result, r.error
	case <-timer.C:
		logging.LogError(c.FormatLogMessage("error", "timed out during telnet authentication"))

		return []byte{}, ErrAuthTimeout
	}
}
