//go:build !windows
// +build !windows

package transport

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

const (
	// SystemTransport is the default "system" (/bin/ssh wrapper) transport for scrapligo.
	SystemTransport = "system"

	defaultOpenBin = "ssh"
)

// NewSystemTransport returns an instance of System transport.
func NewSystemTransport(a *SSHArgs) (*System, error) {
	t := &System{
		SSHArgs:  a,
		openBin:  defaultOpenBin,
		openArgs: make([]string, 0),
		fd:       nil,
	}

	return t, nil
}

// System is the default (/bin/ssh wrapper) transport object.
type System struct {
	SSHArgs   *SSHArgs
	ExtraArgs []string
	openBin   string
	openArgs  []string
	fd        *os.File
}

func (t *System) buildOpenArgs(a *Args) {
	if len(t.openArgs) > 0 {
		t.openArgs = []string{}
	}

	t.openArgs = []string{
		a.Host,
		"-p",
		fmt.Sprintf("%d", a.Port),
		"-o",
		fmt.Sprintf("ConnectTimeout=%d", int(a.TimeoutSocket.Seconds())),
		"-o",
		fmt.Sprintf("ServerAliveInterval=%d", int(a.TimeoutSocket.Seconds())),
	}

	if a.User != "" {
		t.openArgs = append(
			t.openArgs,
			"-l",
			a.User,
		)
	}

	if t.SSHArgs.StrictKey {
		t.openArgs = append(
			t.openArgs,
			"-o",
			"StrictHostKeyChecking=yes",
		)

		if t.SSHArgs.KnownHostsFile != "" {
			t.openArgs = append(
				t.openArgs,
				"-o",
				fmt.Sprintf("UserKnownHostsFile=%s", t.SSHArgs.KnownHostsFile),
			)
		}
	} else {
		t.openArgs = append(
			t.openArgs,
			"-o",
			"StrictHostKeyChecking=no",
			"-o",
			"UserKnownHostsFile=/dev/null",
		)
	}

	if t.SSHArgs.ConfigFile != "" {
		t.openArgs = append(
			t.openArgs,
			"-F",
			t.SSHArgs.ConfigFile,
		)
	} else {
		t.openArgs = append(
			t.openArgs,
			"-F",
			"/dev/null",
		)
	}

	if len(t.ExtraArgs) > 0 {
		t.openArgs = append(
			t.openArgs,
			t.ExtraArgs...,
		)
	}
}

func (t *System) open(a *Args) error {
	if len(t.openArgs) == 0 {
		t.buildOpenArgs(a)
	}

	a.l.Debugf("opening system transport with bin '%s' and args '%s'", t.openBin, t.openArgs)

	c := exec.Command(t.openBin, t.openArgs...) //nolint:gosec

	var err error

	t.fd, err = pty.StartWithSize(
		c,
		&pty.Winsize{
			Rows: uint16(a.TermHeight),
			Cols: uint16(a.TermWidth),
		},
	)
	if err != nil {
		a.l.Criticalf("encountered error spawning pty, error: %s", err)

		return err
	}

	return nil
}

func (t *System) openNetconf(a *Args) error {
	if len(t.openArgs) == 0 {
		t.buildOpenArgs(a)
	}

	t.openArgs = append(t.openArgs, "-s", "netconf")

	a.l.Debugf("opening system transport with bin '%s' and args '%s'", t.openBin, t.openArgs)

	c := exec.Command(t.openBin, t.openArgs...) //nolint:gosec

	var err error

	t.fd, err = pty.Start(c)

	if err != nil {
		a.l.Criticalf("encountered error spawning pty, error: %s", err)

		return err
	}

	return nil
}

// Open opens the System transport.
func (t *System) Open(a *Args) error {
	if t.SSHArgs.NetconfConnection {
		return t.openNetconf(a)
	}

	return t.open(a)
}

// Close closes the System transport.
func (t *System) Close() error {
	err := t.fd.Close()

	t.fd = nil

	return err
}

// IsAlive returns true if the System transport file descriptor is not nil.
func (t *System) IsAlive() bool {
	return t.fd != nil
}

// Read reads n bytes from the transport.
func (t *System) Read(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := t.fd.Read(b)

	return b, err
}

// Write writes bytes b to the transport.
func (t *System) Write(b []byte) error {
	_, err := t.fd.Write(b)

	return err
}
