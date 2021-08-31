//go:build !windows
// +build !windows

package transport

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/scrapli/scrapligo/logging"
)

const sshCmd = "ssh"

// System the "system" (pty subprocess wrapper) transport option for scrapligo.
type System struct {
	SystemTransportArgs *SystemTransportArgs
	fileObj             *os.File
	OpenCmd             []string
	ExecCmd             string
}

// SystemTransport interface describes system transport specific methods.
type SystemTransport interface {
	SetOpenCmd([]string)
	SetExecCmd(string)
}

// SystemTransportArgs struct representing attributes required for the System transport.
type SystemTransportArgs struct {
	AuthPrivateKey    string
	AuthStrictKey     bool
	SSHConfigFile     string
	SSHKnownHostsFile string
	NetconfForcePty   *bool
}

// SetOpenCmd sets the open command string slice; arguments used for opening the connection.
func (t *System) SetOpenCmd(openCmd []string) {
	t.OpenCmd = openCmd
}

// SetExecCmd sets the exec command string, binary used for opening the connection.
func (t *System) SetExecCmd(execCmd string) {
	t.ExecCmd = execCmd
}

func (t *System) buildOpenCmd(baseArgs *BaseTransportArgs) {
	// base open command arguments; the exec command itself will be passed in open()
	// need to add user arguments could go here at some point
	t.OpenCmd = append(
		t.OpenCmd,
		baseArgs.Host,
		"-p",
		fmt.Sprintf("%d", baseArgs.Port),
		"-o",
		fmt.Sprintf("ConnectTimeout=%d", int(baseArgs.TimeoutSocket.Seconds())),
		"-o",
		fmt.Sprintf("ServerAliveInterval=%d", int(baseArgs.TimeoutTransport.Seconds())),
	)

	if t.SystemTransportArgs.AuthPrivateKey != "" {
		t.OpenCmd = append(
			t.OpenCmd,
			"-i",
			t.SystemTransportArgs.AuthPrivateKey,
		)
	}

	if baseArgs.AuthUsername != "" {
		t.OpenCmd = append(
			t.OpenCmd,
			"-l",
			baseArgs.AuthUsername,
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

func (t *System) Open(baseArgs *BaseTransportArgs) error {
	if t.OpenCmd == nil {
		t.buildOpenCmd(baseArgs)
	}

	if t.ExecCmd == "" {
		t.ExecCmd = sshCmd
	}

	logging.LogDebug(
		FormatLogMessage(baseArgs,
			"debug",
			fmt.Sprintf(
				"\"attempting to open transport connection with the following command: %s",
				t.OpenCmd,
			),
		),
	)

	command := exec.Command(t.ExecCmd, t.OpenCmd...) //nolint:gosec
	fileObj, err := pty.StartWithSize(
		command,
		&pty.Winsize{
			Rows: uint16(baseArgs.PtyHeight),
			Cols: uint16(baseArgs.PtyWidth),
		},
	)

	if err == nil {
		t.fileObj = fileObj
	}

	return err
}

func (t *System) OpenNetconf(baseArgs *BaseTransportArgs) error {
	if t.OpenCmd == nil {
		t.buildOpenCmd(baseArgs)

		if t.SystemTransportArgs.NetconfForcePty == nil || *t.SystemTransportArgs.NetconfForcePty {
			t.OpenCmd = append(t.OpenCmd, "-tt")
		}

		t.OpenCmd = append(t.OpenCmd,
			"-s",
			"netconf",
		)
	}

	if t.ExecCmd == "" {
		t.ExecCmd = sshCmd
	}

	logging.LogDebug(
		FormatLogMessage(baseArgs,
			"debug",
			fmt.Sprintf(
				"\"attempting to open netconf transport connection with the following command: %s",
				t.OpenCmd,
			),
		),
	)

	command := exec.Command(t.ExecCmd, t.OpenCmd...) //nolint:gosec
	fileObj, err := pty.Start(command)

	if err == nil {
		t.fileObj = fileObj
	}

	return err
}

func (t *System) Close() error {
	err := t.fileObj.Close()
	t.fileObj = nil

	return err
}

func (t *System) IsAlive() bool {
	return t.fileObj != nil
}

func (t *System) Read(n int) *ReadResult {
	b := make([]byte, n)
	_, err := t.fileObj.Read(b)

	if err != nil {
		return &ReadResult{
			Result: nil,
			Error:  ErrTransportFailure,
		}
	}

	return &ReadResult{
		Result: b,
		Error:  nil,
	}
}

func (t *System) Write(channelInput []byte) error {
	_, err := t.fileObj.Write(channelInput)

	return err
}
