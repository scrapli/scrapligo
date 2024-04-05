package transport_test

import (
	"flag"
	"net"
	"os"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/scrapli/scrapligo/util"
)

var (
	update = flag.Bool( //nolint
		"update",
		false,
		"update the golden files",
	)
	functional = flag.Bool( //nolint
		"functional",
		false,
		"execute functional tests",
	)
	platforms = flag.String( //nolint
		"platforms",
		util.All,
		"comma sep list of platform(s) to target",
	)
	transports = flag.String( //nolint
		"transports",
		util.All,
		"comma sep list of transport(s) to target",
	)
)

func spawnDumbSever(
	t *testing.T,
	handlerDone chan struct{},
	handlerErrChan chan error,
	writeOuts ...string,
) {
	t.Helper()

	serverConfig := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	privateKeyBytes, err := os.ReadFile("test-fixtures/dumbserver")
	if err != nil {
		handlerErrChan <- err

		return
	}

	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		handlerErrChan <- err

		return
	}

	serverConfig.AddHostKey(privateKey)

	listener, err := net.Listen("tcp", "0.0.0.0:2222") //nolint: gosec
	if err != nil {
		handlerErrChan <- err

		return
	}

	defer listener.Close() //nolint:errcheck

	// not wrapping this in goroutine or handling multiple connections because we dont care
	tcpConn, acceptErr := listener.Accept()
	if acceptErr != nil {
		handlerErrChan <- err

		return
	}

	defer tcpConn.Close() //nolint:errcheck

	_, sshChans, reqs, connErr := ssh.NewServerConn(tcpConn, serverConfig)
	if connErr != nil {
		handlerErrChan <- err

		return
	}

	go ssh.DiscardRequests(reqs)

	// we only want the first/only channel for the test, nice and dumb and simple
	sshChan := <-sshChans

	time.Sleep(100 * time.Millisecond)

	sshConn, _, chanErr := sshChan.Accept()
	if chanErr != nil {
		handlerErrChan <- err

		return
	}

	for _, writeOut := range writeOuts {
		_, err = sshConn.Write([]byte(writeOut))
		if err != nil {
			handlerErrChan <- err

			return
		}
	}

	close(handlerDone)

	// block here to hold the connection open (meaning dont let the defer'd calls run)
	time.Sleep(60 * time.Second)
}
