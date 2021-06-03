package transport

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/scrapli/scrapligo/logging"

	"golang.org/x/crypto/ssh"
)

// ErrNotImplemented error raised when trying to use netconf w/ standard transport.
var ErrNotImplemented = errors.New("not implemented")

// Standard the "standard" (standard library) transport option for scrapligo.
type Standard struct {
	BaseTransportArgs     *BaseTransportArgs
	StandardTransportArgs *StandardTransportArgs
	client                *ssh.Client
	session               *ssh.Session
	writer                io.WriteCloser
	reader                io.Reader
}

// StandardTransportArgs struct representing attributes required for the Standard transport.
type StandardTransportArgs struct {
	AuthPassword      string
	AuthPrivateKey    string
	AuthStrictKey     bool
	SSHConfigFile     string
	SSHKnownHostsFile string
}

func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(
		k.Marshal(),
	) // e.g. "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTY...."
}

// https://stackoverflow.com/questions/44269142/ \
// golang-ssh-getting-must-specify-hoskeycallback-error-despite-setting-it-to-n
// basically need to parse ssh config like scrapli does... at some point.
func trustedHostKeyCallback(trustedKey string) ssh.HostKeyCallback {
	if trustedKey == "" {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			log.Printf(
				"ssh key verification is *NOT* in effect: to fix, add this trustedKey: %q",
				keyString(k),
			)

			return nil
		}
	}

	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := keyString(k)
		if trustedKey != ks {
			return ErrKeyVerificationFailed
		}

		return nil
	}
}

func (t *Standard) open(cfg *ssh.ClientConfig) error {
	var err error

	t.client, err = ssh.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", t.BaseTransportArgs.Host, t.BaseTransportArgs.Port),
		cfg,
	)
	if err != nil {
		logging.ErrorLog(
			t.FormatLogMessage("error", fmt.Sprintf("error connecting to host: %v", err)),
		)

		return err
	}

	t.session, err = t.client.NewSession()
	if err != nil {
		logging.ErrorLog(
			t.FormatLogMessage("error", fmt.Sprintf("error allocating session: %v", err)),
		)

		return err
	}

	t.writer, err = t.session.StdinPipe()
	if err != nil {
		logging.ErrorLog(
			t.FormatLogMessage("error", fmt.Sprintf("error allocating writer: %v", err)),
		)

		return err
	}

	t.reader, err = t.session.StdoutPipe()
	if err != nil {
		logging.ErrorLog(
			t.FormatLogMessage("error", fmt.Sprintf("error allocating reader: %v", err)),
		)

		return err
	}

	return nil
}

func (t *Standard) openBase() error {
	/* #nosec G106 */
	hostKeyCallback := ssh.InsecureIgnoreHostKey()
	if t.StandardTransportArgs.AuthStrictKey {
		// trustedKey will need to be gleaned from known hosts how scrapli does at some point
		hostKeyCallback = trustedHostKeyCallback("")
	}

	authMethods := make([]ssh.AuthMethod, 0)

	if t.StandardTransportArgs.AuthPrivateKey != "" {
		key, err := ioutil.ReadFile(t.StandardTransportArgs.AuthPrivateKey)
		if err != nil {
			return err
		}

		signer, err := ssh.ParsePrivateKey(key)

		if err != nil {
			logging.ErrorLog(
				t.FormatLogMessage("error", fmt.Sprintf("unable to parse private key: %v", err)),
			)

			return err
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if t.StandardTransportArgs.AuthPassword != "" {
		authMethods = append(authMethods, ssh.Password(t.StandardTransportArgs.AuthPassword),
			ssh.KeyboardInteractive(
				func(user, instruction string, questions []string, echos []bool) ([]string, error) {
					answers := make([]string, len(questions))
					for i := range answers {
						answers[i] = t.StandardTransportArgs.AuthPassword
					}

					return answers, nil
				},
			))
	}

	cfg := &ssh.ClientConfig{
		User:            t.BaseTransportArgs.AuthUsername,
		Auth:            authMethods,
		Timeout:         *t.BaseTransportArgs.TimeoutSocket,
		HostKeyCallback: hostKeyCallback,
	}

	err := t.open(cfg)
	if err != nil {
		return err
	}

	// not sure what to do about the tty speeds... figured lets just go fast?
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	err = t.session.RequestPty(
		"xterm",
		t.BaseTransportArgs.PtyHeight,
		t.BaseTransportArgs.PtyWidth,
		modes,
	)
	if err != nil {
		return err
	}

	return nil
}

// Open open a standard ssh connection.
func (t *Standard) Open() error {
	err := t.openBase()
	if err != nil {
		return err
	}

	err = t.session.Shell()
	if err != nil {
		return err
	}

	return nil
}

// OpenNetconf open a netconf connection.
func (t *Standard) OpenNetconf() error {
	err := t.openBase()
	if err != nil {
		logging.LogError(
			fmt.Sprintf(
				"failed opening base connection, cant attempt to open netconf connection; error: %v",
				err,
			),
		)

		return err
	}

	err = t.session.RequestSubsystem("netconf")
	if err != nil {
		logging.LogError(fmt.Sprintf("failed opening netconf subsystem; error: %v", err))
		return err
	}

	return nil
}

// Close close the transport connection to the device.
func (t *Standard) Close() error {
	err := t.session.Close()
	t.session = nil

	logging.LogDebug(t.FormatLogMessage("debug", "transport connection to host closed"))

	return err
}

func (t *Standard) read() *transportResult {
	b := make([]byte, 0, ReadSize)
	_, err := t.reader.Read(b)

	if err != nil {
		return &transportResult{
			result: nil,
			error:  ErrTransportFailure,
		}
	}

	return &transportResult{
		result: b,
		error:  nil,
	}
}

// Read read bytes from the transport.
func (t *Standard) Read() ([]byte, error) {
	b, err := transportTimeout(
		*t.BaseTransportArgs.TimeoutTransport,
		t.read,
	)

	if err != nil {
		logging.LogError(t.FormatLogMessage("error", "timed out reading from transport"))
		return b, err
	}

	return b, nil
}

// Write write bytes to the transport.
func (t *Standard) Write(channelInput []byte) error {
	_, err := t.writer.Write(channelInput)
	if err != nil {
		return err
	}

	return nil
}

// IsAlive indicate if the transport is alive or not.
func (t *Standard) IsAlive() bool {
	return t.session != nil
}

// FormatLogMessage formats log message payload, adding contextual info about the host.
func (t *Standard) FormatLogMessage(level, msg string) string {
	return logging.FormatLogMessage(level, t.BaseTransportArgs.Host, t.BaseTransportArgs.Port, msg)
}
