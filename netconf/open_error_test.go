package netconf_test

import (
	"context"
	"errors"
	"testing"
	"time"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
	scrapligonetconf "github.com/scrapli/scrapligo/v2/netconf"
	scrapligooptions "github.com/scrapli/scrapligo/v2/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/v2/testhelper"
)

func TestOpenErrors(t *testing.T) {
	cases := map[string]struct {
		port     func(t *testing.T) uint16
		options  []scrapligooptions.Option
		ctx      func(t *testing.T) context.Context
		expected error
	}{
		"connection-refused": {
			port:     scrapligotesthelper.ClosedPort,
			expected: scrapligoerrors.ErrSession,
		},
		"operation-timeout": {
			port: scrapligotesthelper.SilentServer,
			options: []scrapligooptions.Option{
				scrapligooptions.WithOperationTimeout(time.Second),
			},
			expected: scrapligoerrors.ErrTimeout,
		},
		"context-deadline": {
			port: scrapligotesthelper.SilentServer,
			ctx: func(t *testing.T) context.Context {
				t.Helper()

				ctx, cancel := context.WithTimeout(t.Context(), time.Second)
				t.Cleanup(cancel)

				return ctx
			},
			expected: context.DeadlineExceeded,
		},
		"context-cancelled": {
			port: scrapligotesthelper.SilentServer,
			ctx: func(t *testing.T) context.Context {
				t.Helper()

				ctx, cancel := context.WithCancel(t.Context())
				time.AfterFunc(time.Second, cancel)
				t.Cleanup(cancel)

				return ctx
			},
			expected: context.Canceled,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			if c.ctx != nil {
				ctx = c.ctx(t)
			}

			d, err := scrapligonetconf.NewNetconf(
				"127.0.0.1",
				append(c.options, scrapligooptions.WithPort(c.port(t)))...,
			)
			if err != nil {
				t.Fatal(err)
			}

			_, err = d.Open(ctx)
			if !errors.Is(err, c.expected) {
				t.Fatalf("expected errors.Is(err, %v), got %v", c.expected, err)
			}
		})
	}
}
