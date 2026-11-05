package transport_test

import (
	"bytes"
	"errors"
	"net"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/logging"
	"github.com/scrapli/scrapligo/transport"
)

func TestStandardTransportHostKeyError(t *testing.T) {
	privateKeyBytes, err := os.ReadFile("test-fixtures/dumbserver")
	if err != nil {
		t.Fatal(err)
	}

	hostKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	serverConfig := &ssh.ServerConfig{
		NoClientAuth: true,
	}
	serverConfig.AddHostKey(hostKey)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close() //nolint:errcheck

	go func() {
		tcpConn, acceptErr := listener.Accept()
		if acceptErr != nil {
			return
		}

		defer tcpConn.Close() //nolint:errcheck

		_, _, _, _ = ssh.NewServerConn(tcpConn, serverConfig)
	}()

	knownHostsFile := filepath.Join(t.TempDir(), "known_hosts")

	err = os.WriteFile(knownHostsFile, nil, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	sshArgs, err := transport.NewSSHArgs(options.WithSSHKnownHostsFile(knownHostsFile))
	if err != nil {
		t.Fatal(err)
	}

	tp, err := transport.NewStandardTransport(sshArgs)
	if err != nil {
		t.Fatal(err)
	}

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("listener address is not a TCP address")
	}

	openArgs, err := transport.NewArgs(
		&logging.Instance{},
		"127.0.0.1",
		options.WithPort(addr.Port),
		options.WithAuthUsername("whatever"),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = tp.Open(openArgs)

	var hostKeyErr *transport.HostKeyError
	if !errors.As(err, &hostKeyErr) {
		t.Fatalf("expected HostKeyError, got: %v", err)
	}

	if !bytes.Equal(hostKeyErr.HostKey.Marshal(), hostKey.PublicKey().Marshal()) {
		t.Fatal("HostKeyError does not hold the host key presented by the server")
	}

	var keyErr *knownhosts.KeyError
	if !errors.As(err, &keyErr) {
		t.Fatalf("expected HostKeyError to wrap knownhosts.KeyError, got: %v", err)
	}
}
