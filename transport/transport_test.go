package transport_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/gliderlabs/ssh"

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

func dummyDevice(
	t *testing.T,
	port int,
	sshHandler ssh.Handler,
	netconfHandler ssh.SubsystemHandler,
) *ssh.Server {
	s := &ssh.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	if sshHandler != nil {
		s.Handler = sshHandler
	}

	if netconfHandler != nil {
		s.SubsystemHandlers = map[string]ssh.SubsystemHandler{
			"netconf": netconfHandler,
		}
	}

	// run the server
	go func() {
		if err := s.ListenAndServe(); err != nil {
			t.Log(err)
		}
	}()

	return s
}
