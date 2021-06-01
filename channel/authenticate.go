package channel

import (
	"bytes"
	"time"

	"github.com/scrapli/scrapligo/logging"
)

func (c *Channel) authenticateSSH(authPassword, authPassphrase []byte) *channelResult {
	logging.LogDebug(c.FormatLogMessage("debug", "attempting in channel ssh authentication"))

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

		if passwordMatch {
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

// AuthenticateSSH Handle "in channel" SSH authentication (for "system" transport).
func (c *Channel) AuthenticateSSH(authPassword, authPassphrase string) ([]byte, error) {
	var _c = make(chan *channelResult, 1)

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
		logging.LogError(c.FormatLogMessage("error", "timed out during authentication"))

		return []byte{}, ErrAuthTimeout
	}
}
