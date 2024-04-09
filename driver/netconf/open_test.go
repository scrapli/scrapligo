package netconf_test

import (
	"testing"
	"time"

	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/transport"
)

type testCase struct {
	Description string
	PayloadFile string

	ExpectErr bool
}

func testOpen(testName string, testCase *testCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		doneChan := make(chan struct{})
		go func() {
			d, err := netconf.NewDriver(
				"dummy",
				options.WithTransportType(transport.FileTransport),
				options.WithFileTransportFile(resolveFile(t, testCase.PayloadFile)),
				// options.WithTransportReadSize(1),
				options.WithReadDelay(0),
				options.WithTimeoutOps(1*time.Second),
			)
			if err != nil {
				if testCase.ExpectErr {
					t.Logf("%s: encountered error creating network Driver, error: %s", testName, err)
					} else {
					t.Fatalf("%s: encountered error creating network Driver, error: %s", testName, err)
				}
			}

			err = d.Open()
			if err != nil {
				if testCase.ExpectErr {
					t.Logf("%s: encountered error opening netconf Driver, error: %s", testName, err)
					} else {
					t.Fatalf("%s: encountered error opening netconf Driver, error: %s", testName, err)
				}
			}

			defer func() {
				if err := d.Close(); err != nil {
					t.Logf("failed to close driver: %s", err)
				}
				close(doneChan)
			}()
		}()

		select {
		case <- doneChan:
		case <- time.After(5*time.Second):
			t.Error("timeout waiting for open to complete")
		}
	}
}

func TestOpen(t *testing.T) {
	cases := map[string]*testCase{
		"server-capabilities-truncated": {
			Description: "server capabilities truncated",
			PayloadFile: "open-server-capabilities-truncated.txt",
			ExpectErr: true,
		},
		"simple": {
			Description: "simple open",
			PayloadFile: "open-simple.txt",
			ExpectErr: false,
		},
	}

	for testName, testCase := range cases {
		f := testOpen(testName, testCase)
		t.Run(testName, f)
	}
}