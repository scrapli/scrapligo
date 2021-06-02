// +build !windows

package transport

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/scrapli/scrapligo/logging"

	"github.com/creack/pty"
)

// System the "system" (pty subprocess wrapper) transport option for scrapligo.
type System struct {
	BaseTransportArgs   *BaseTransportArgs
	SystemTransportArgs *SystemTransportArgs
	sessionFd           *os.File
	openCmd             []string
}

// SystemTransportArgs struct representing attributes required for the System transport.
type SystemTransportArgs struct {
	AuthPrivateKey    string
	AuthStrictKey     bool
	SSHConfigFile     string
	SSHKnownHostsFile string
}

func (t *System) buildOpenCmd() {
	// base ssh arguments; "ssh" itself passed in Open()
	// need to add user arguments could go here at some point
	t.openCmd = append(
		t.openCmd,
		t.BaseTransportArgs.Host,
		"-p",
		fmt.Sprintf("%d", t.BaseTransportArgs.Port),
		"-o",
		fmt.Sprintf("ConnectTimeout=%d", int(t.BaseTransportArgs.TimeoutSocket.Seconds())),
		"-o",
		fmt.Sprintf("ServerAliveInterval=%d", int(t.BaseTransportArgs.TimeoutTransport.Seconds())),
	)

	if t.SystemTransportArgs.AuthPrivateKey != "" {
		t.openCmd = append(
			t.openCmd,
			"-i",
			t.SystemTransportArgs.AuthPrivateKey,
		)
	}

	if t.BaseTransportArgs.AuthUsername != "" {
		t.openCmd = append(
			t.openCmd,
			"-l",
			t.BaseTransportArgs.AuthUsername,
		)
	}

	if !t.SystemTransportArgs.AuthStrictKey {
		t.openCmd = append(
			t.openCmd,
			"-o",
			"StrictHostKeyChecking=no",
			"-o",
			"UserKnownHostsFile=/dev/null",
		)
	} else {
		t.openCmd = append(
			t.openCmd,
			"-o",
			"StrictHostKeyChecking=yes",
		)

		if t.SystemTransportArgs.SSHKnownHostsFile != "" {
			t.openCmd = append(
				t.openCmd,
				"-o",
				fmt.Sprintf("UserKnownHostsFile=%s", t.SystemTransportArgs.SSHKnownHostsFile),
			)
		}
	}

	if t.SystemTransportArgs.SSHConfigFile != "" {
		t.openCmd = append(
			t.openCmd,
			"-F",
			t.SystemTransportArgs.SSHConfigFile,
		)
	} else {
		t.openCmd = append(
			t.openCmd,
			"-F",
			"/dev/null",
		)
	}
}

// Open open a standard ssh connection.
func (t *System) Open() error {
	t.buildOpenCmd()

	logging.LogDebug(
		t.FormatLogMessage(
			"debug",
			fmt.Sprintf(
				"\"attempting to open transport connection with the following command: %s",
				t.openCmd,
			),
		),
	)

	sshCommand := exec.Command("ssh", t.openCmd...)
	sessionFd, err := pty.StartWithSize(
		sshCommand,
		&pty.Winsize{
			Rows: uint16(t.BaseTransportArgs.PtyHeight),
			Cols: uint16(t.BaseTransportArgs.PtyWidth),
		},
	)

	if err != nil {
		logging.ErrorLog(t.FormatLogMessage("error", "failed opening transport connection to host"))

		return err
	}

	logging.LogDebug(t.FormatLogMessage("debug", "transport connection to host opened"))

	t.sessionFd = sessionFd

	return err
}

// OpenNetconf open a netconf connection.
func (t *System) OpenNetconf() error {
	t.buildOpenCmd()

	t.openCmd = append(t.openCmd,
		"-tt",
		"-s",
		"netconf",
	)

	logging.LogDebug(
		t.FormatLogMessage(
			"debug",
			fmt.Sprintf(
				"\"attempting to open netconf transport connection with the following command: %s",
				t.openCmd,
			),
		),
	)

	sshCommand := exec.Command("ssh", t.openCmd...)
	sessionFd, err := pty.Start(sshCommand)

	if err != nil {
		logging.ErrorLog(
			t.FormatLogMessage("error", "failed opening netconf transport connection to host"),
		)

		return err
	}

	logging.LogDebug(t.FormatLogMessage("debug", "netconf transport connection to host opened"))

	t.sessionFd = sessionFd

	return err
}

// Close close the transport connection to the device.
func (t *System) Close() error {
	err := t.sessionFd.Close()
	t.sessionFd = nil
	logging.LogDebug(t.FormatLogMessage("debug", "transport connection to host closed"))

	return err
}

func (t *System) read() *transportResult {
	b := make([]byte, ReadSize)
	_, err := t.sessionFd.Read(b)

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
func (t *System) Read() ([]byte, error) {
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
func (t *System) Write(channelInput []byte) error {
	_, err := t.sessionFd.Write(channelInput)
	if err != nil {
		return err
	}

	return nil
}

// IsAlive indicate if the transport is alive or not.
func (t *System) IsAlive() bool {
	return t.sessionFd != nil
}

// FormatLogMessage formats log message payload, adding contextual info about the host.
func (t *System) FormatLogMessage(level, msg string) string {
	return logging.FormatLogMessage(level, t.BaseTransportArgs.Host, t.BaseTransportArgs.Port, msg)
}
