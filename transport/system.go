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
	fileObj             *os.File
	OpenCmd             []string
	ExecCmd             string
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
	t.OpenCmd = append(
		t.OpenCmd,
		t.BaseTransportArgs.Host,
		"-p",
		fmt.Sprintf("%d", t.BaseTransportArgs.Port),
		"-o",
		fmt.Sprintf("ConnectTimeout=%d", int(t.BaseTransportArgs.TimeoutSocket.Seconds())),
		"-o",
		fmt.Sprintf("ServerAliveInterval=%d", int(t.BaseTransportArgs.TimeoutTransport.Seconds())),
	)

	if t.SystemTransportArgs.AuthPrivateKey != "" {
		t.OpenCmd = append(
			t.OpenCmd,
			"-i",
			t.SystemTransportArgs.AuthPrivateKey,
		)
	}

	if t.BaseTransportArgs.AuthUsername != "" {
		t.OpenCmd = append(
			t.OpenCmd,
			"-l",
			t.BaseTransportArgs.AuthUsername,
		)
	}

	if !t.SystemTransportArgs.AuthStrictKey {
		t.OpenCmd = append(
			t.OpenCmd,
			"-o",
			"StrictHostKeyChecking=no",
			"-o",
			"UserKnownHostsFile=/dev/null",
		)
	} else {
		t.OpenCmd = append(
			t.OpenCmd,
			"-o",
			"StrictHostKeyChecking=yes",
		)

		if t.SystemTransportArgs.SSHKnownHostsFile != "" {
			t.OpenCmd = append(
				t.OpenCmd,
				"-o",
				fmt.Sprintf("UserKnownHostsFile=%s", t.SystemTransportArgs.SSHKnownHostsFile),
			)
		}
	}

	if t.SystemTransportArgs.SSHConfigFile != "" {
		t.OpenCmd = append(
			t.OpenCmd,
			"-F",
			t.SystemTransportArgs.SSHConfigFile,
		)
	} else {
		t.OpenCmd = append(
			t.OpenCmd,
			"-F",
			"/dev/null",
		)
	}
}

// Open open a standard ssh connection.
func (t *System) Open() error {
	if t.OpenCmd == nil {
		t.buildOpenCmd()
	}

	if t.ExecCmd == "" {
		t.ExecCmd = "ssh"
	}

	logging.LogDebug(
		t.FormatLogMessage(
			"debug",
			fmt.Sprintf(
				"\"attempting to open transport connection with the following command: %s",
				t.OpenCmd,
			),
		),
	)

	sshCommand := exec.Command(t.ExecCmd, t.OpenCmd...)
	fileObj, err := pty.StartWithSize(
		sshCommand,
		&pty.Winsize{
			Rows: uint16(t.BaseTransportArgs.PtyHeight),
			Cols: uint16(t.BaseTransportArgs.PtyWidth),
		},
	)

	if err != nil {
		logging.LogError(t.FormatLogMessage("error", "failed opening transport connection to host"))

		return err
	}

	logging.LogDebug(t.FormatLogMessage("debug", "transport connection to host opened"))

	t.fileObj = fileObj

	return err
}

// OpenNetconf open a netconf connection.
func (t *System) OpenNetconf() error {
	t.buildOpenCmd()

	t.OpenCmd = append(t.OpenCmd,
		"-tt",
		"-s",
		"netconf",
	)

	logging.LogDebug(
		t.FormatLogMessage(
			"debug",
			fmt.Sprintf(
				"\"attempting to open netconf transport connection with the following command: %s",
				t.OpenCmd,
			),
		),
	)

	sshCommand := exec.Command("ssh", t.OpenCmd...)
	fileObj, err := pty.Start(sshCommand)

	if err != nil {
		logging.LogError(
			t.FormatLogMessage("error", "failed opening netconf transport connection to host"),
		)

		return err
	}

	logging.LogDebug(t.FormatLogMessage("debug", "netconf transport connection to host opened"))

	t.fileObj = fileObj

	return err
}

// Close close the transport connection to the device.
func (t *System) Close() error {
	err := t.fileObj.Close()
	t.fileObj = nil
	logging.LogDebug(t.FormatLogMessage("debug", "transport connection to host closed"))

	return err
}

func (t *System) read(n int) *transportResult {
	b := make([]byte, n)
	_, err := t.fileObj.Read(b)

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
		ReadSize,
	)

	if err != nil {
		logging.LogError(t.FormatLogMessage("error", "timed out reading from transport"))
		return b, err
	}

	return b, nil
}

// ReadN read N bytes from the transport.
func (t *System) ReadN(n int) ([]byte, error) {
	b, err := transportTimeout(
		*t.BaseTransportArgs.TimeoutTransport,
		t.read,
		n,
	)

	if err != nil {
		logging.LogError(t.FormatLogMessage("error", "timed out reading from transport"))
		return b, err
	}

	return b, nil
}

// Write write bytes to the transport.
func (t *System) Write(channelInput []byte) error {
	_, err := t.fileObj.Write(channelInput)
	if err != nil {
		return err
	}

	return nil
}

// IsAlive indicate if the transport is alive or not.
func (t *System) IsAlive() bool {
	return t.fileObj != nil
}

// FormatLogMessage formats log message payload, adding contextual info about the host.
func (t *System) FormatLogMessage(level, msg string) string {
	return logging.FormatLogMessage(level, t.BaseTransportArgs.Host, t.BaseTransportArgs.Port, msg)
}
