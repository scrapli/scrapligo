package channel

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/scrapli/scrapligo/util"
)

const (
	usernameSeenMax   = 2
	passwordSeenMax   = 2
	passphraseSeenMax = 2
)

type authPatterns struct {
	username   *regexp.Regexp
	password   *regexp.Regexp
	passphrase *regexp.Regexp
}

var (
	authPatternsInstance     *authPatterns //nolint:gochecknoglobals
	authPatternsInstanceOnce sync.Once     //nolint:gochecknoglobals
)

func getAuthPatterns() *authPatterns {
	authPatternsInstanceOnce.Do(func() {
		authPatternsInstance = &authPatterns{
			username:   regexp.MustCompile(`(?im)^(.*username:)|(.*login:)\s?$`),
			password:   regexp.MustCompile(`(?im)(.*@.*)?password:\s?$`),
			passphrase: regexp.MustCompile(`(?i)enter passphrase for key`),
		}
	})

	return authPatternsInstance
}

func (c *Channel) authenticateSSH(p, pp []byte) *result {
	pCount := 0

	ppCount := 0

	var b []byte

	for {
		nb, err := c.ReadUntilAnyPrompt(
			[]*regexp.Regexp{c.PromptPattern, c.PasswordPattern, c.PassphrasePattern},
		)
		if err != nil {
			return &result{nil, err}
		}

		b = append(b, nb...)

		if c.PromptPattern.Match(b) {
			return &result{b, nil}
		}

		if c.PasswordPattern.Match(b) { //nolint:nestif
			b = []byte{}

			pCount++

			if pCount > passwordSeenMax {
				c.l.Critical("password prompt seen multiple times, assuming authentication failed")

				return &result{
					nil,
					fmt.Errorf(
						"%w: password prompt seen multiple times, assuming authentication failed",
						util.ErrAuthError,
					),
				}
			}

			err = c.WriteAndReturn(p, true)
			if err != nil {
				return &result{nil, err}
			}
		} else if c.PassphrasePattern.Match(b) {
			b = []byte{}

			ppCount++

			if ppCount > passphraseSeenMax {
				c.l.Critical(
					"private key passphrase prompt seen multiple times," +
						" assuming authentication failed",
				)

				return &result{
					nil,
					fmt.Errorf(
						"%w: private key passphrase prompt seen multiple times,"+
							" assuming authentication failed",
						util.ErrAuthError,
					),
				}
			}

			err = c.WriteAndReturn(pp, true)
			if err != nil {
				return &result{nil, err}
			}
		}
	}
}

// AuthenticateSSH handles "in channel" SSH authentication.
func (c *Channel) AuthenticateSSH(p, pp []byte) ([]byte, error) {
	cr := make(chan *result)

	go func() {
		cr <- c.authenticateSSH(p, pp)
	}()

	t := time.NewTimer(c.TimeoutOps)

	select {
	case r := <-cr:
		return r.b, r.err
	case <-t.C:
		c.l.Critical("channel timeout during in channel ssh authentication")

		return nil, fmt.Errorf(
			"%w: channel timeout during in channel ssh authentication",
			util.ErrTimeoutError,
		)
	}
}

func (c *Channel) authenticateTelnet(u, p []byte) *result {
	uCount := 0

	pCount := 0

	var b []byte

	for {
		nb, err := c.ReadUntilAnyPrompt(
			[]*regexp.Regexp{c.PromptPattern, c.UsernamePattern, c.PasswordPattern},
		)
		if err != nil {
			return &result{nil, err}
		}

		b = append(b, nb...)

		if c.PromptPattern.Match(b) {
			return &result{b, nil}
		}

		if c.UsernamePattern.Match(b) { //nolint:nestif
			b = []byte{}

			uCount++

			if uCount > usernameSeenMax {
				c.l.Critical(
					"username prompt seen multiple times, assuming authentication failed",
				)

				return &result{
					nil,
					fmt.Errorf(
						"%w: username prompt seen multiple times, assuming authentication failed",
						util.ErrAuthError,
					),
				}
			}

			err = c.WriteAndReturn(u, true)
			if err != nil {
				return &result{nil, err}
			}
		} else if c.PasswordPattern.Match(b) {
			b = []byte{}

			pCount++

			if pCount > passwordSeenMax {
				c.l.Critical(
					"password prompt seen multiple times, assuming authentication failed",
				)

				return &result{
					nil,
					fmt.Errorf(
						"%w: password prompt seen multiple times, assuming authentication failed",
						util.ErrAuthError,
					),
				}
			}

			err = c.WriteAndReturn(p, true)
			if err != nil {
				return &result{nil, err}
			}
		}
	}
}

// AuthenticateTelnet handles "in channel" telnet authentication.
func (c *Channel) AuthenticateTelnet(u, p []byte) ([]byte, error) {
	cr := make(chan *result)

	go func() {
		cr <- c.authenticateTelnet(u, p)
	}()

	t := time.NewTimer(c.TimeoutOps)

	select {
	case r := <-cr:
		return r.b, r.err
	case <-t.C:
		c.l.Critical("channel timeout during in channel telnet authentication")

		return nil, fmt.Errorf(
			"%w: channel timeout during in channel telnet authentication",
			util.ErrTimeoutError,
		)
	}
}
