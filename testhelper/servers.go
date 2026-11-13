package testhelper

import (
	"net"
	"sync"
	"testing"
)

func listen(t *testing.T) net.Listener {
	t.Helper()

	lc := net.ListenConfig{}

	l, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	return l
}

func port(t *testing.T, l net.Listener) uint16 {
	t.Helper()

	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener address type %T", l.Addr())
	}

	return uint16(addr.Port) //nolint:gosec
}

// ClosedPort returns a localhost port with nothing listening on it.
func ClosedPort(t *testing.T) uint16 {
	t.Helper()

	l := listen(t)

	p := port(t, l)

	_ = l.Close()

	return p
}

// SilentServer starts a localhost server that accepts connections but never writes anything, and
// returns its port. The server is stopped when the test completes.
func SilentServer(t *testing.T) uint16 {
	t.Helper()

	l := listen(t)

	var (
		connsLock sync.Mutex
		conns     []net.Conn
		wg        sync.WaitGroup
	)

	wg.Go(func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}

			connsLock.Lock()

			conns = append(conns, conn)

			connsLock.Unlock()
		}
	})

	t.Cleanup(func() {
		_ = l.Close()

		wg.Wait()

		for _, conn := range conns {
			_ = conn.Close()
		}
	})

	return port(t, l)
}
