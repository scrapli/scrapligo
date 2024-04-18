package netconf_test

import (
	"testing"
	"time"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/transport"
)

func testOpen(testName string, testCase *util.PayloadTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		doneChan := make(chan struct{})
		errChan := make(chan error)

		go func() {
			defer close(doneChan)

			d, err := netconf.NewDriver(
				"dummy",
				options.WithTransportType(transport.FileTransport),
				options.WithFileTransportFile(resolveFile(t, testCase.PayloadFile)),
				options.WithReadDelay(0),
				options.WithTimeoutOps(1*time.Second),
			)
			if err != nil {
				errChan <- err

				return
			}

			err = d.Open()
			if err != nil {
				if testCase.ExpectErr {
					t.Logf(
						"%s: encountered (expected) error opening netconf Driver, error: %s",
						testName,
						err,
					)
				} else {
					errChan <- err

					return
				}
			}
		}()

		select {
		case <-doneChan:
		case err := <-errChan:
			t.Fatalf("test failed, error: %s", err)
		case <-time.After(5 * time.Second):
			t.Error("timeout waiting for open to complete")
		}
	}
}

func TestOpen(t *testing.T) {
	cases := map[string]*util.PayloadTestCase{
		"server-capabilities-truncated": {
			Description: "server capabilities truncated",
			PayloadFile: "open-server-capabilities-truncated.txt",
			ExpectErr:   true,
		},
		"simple": {
			Description: "simple open",
			PayloadFile: "open-simple.txt",
			ExpectErr:   false,
		},
	}

	for testName, testCase := range cases {
		t.Run(testName, testOpen(testName, testCase))
	}
}
