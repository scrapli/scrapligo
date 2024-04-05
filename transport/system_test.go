package transport_test

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/logging"
	"github.com/scrapli/scrapligo/transport"
)

func TestSystemTransportDontBlockOnClose(t *testing.T) {
	handlerDone := make(chan struct{})
	handlerErrChan := make(chan error)

	go spawnDumbSever(
		t,
		handlerDone,
		handlerErrChan,
		strings.Repeat("z", 8_192),
		strings.Repeat("a", 10),
	)

	sshArgs, err := transport.NewSSHArgs(
		options.WithAuthNoStrictKey(),
		options.WithSSHKnownHostsFile("/dev/null"),
	)
	if err != nil {
		t.Fatal(err)
	}

	tp, err := transport.NewSystemTransport(sshArgs)
	if err != nil {
		t.Fatal(err)
	}

	openArgs, err := transport.NewArgs(
		&logging.Instance{},
		"localhost",
		options.WithPort(2222),
		// doesnt matter dummy server says no auth is ok
		options.WithAuthUsername("whatever"),
		options.WithAuthPassword("whatever"),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = tp.Open(openArgs)
	if err != nil {
		t.Fatal(err)
	}

	doneChan := make(chan struct{})
	errChan := make(chan error)

	go func() {
		defer close(doneChan)

		for {
			_, readErr := tp.Read(81)
			if readErr != nil {
				if readErr == io.EOF {
					return
				}

				errChan <- readErr

				return
			}
		}
	}()

	// wait till the handler is done then we can close the transport to check if we have blocked
	<-handlerDone

	// a small sleep to make sure things have percolated -- without this we get an unrelated error
	// reading from the fd in the transport; we just wanna make sure we dont block so we can ignore
	// that for now
	time.Sleep(time.Second)

	err = tp.Close()
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-doneChan:
		// done chan had a send, so we know the read finished/unblocked
	case handlerErrr := <-handlerErrChan:
		t.Fatalf("error in dumb ssh server, error: %s", handlerErrr)
	case readErr := <-errChan:
		t.Fatalf("non-eof error in read loop, error: %s", readErr)
	case <-time.After(5 * time.Second):
		t.Fatal("read blocked and it should not after transport closure")
	}
}
