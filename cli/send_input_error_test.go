package cli_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/v2/cli"
	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
	scrapligooptions "github.com/scrapli/scrapligo/v2/options"
)

// getUnresponsiveCli returns an opened Cli whose device stops responding after open: the fixture
// holds a recorded open and get prompt, once both have been replayed the test transport has
// nothing more to read.
func getUnresponsiveCli(t *testing.T, opts ...scrapligooptions.Option) *scrapligocli.Cli {
	t.Helper()

	f, err := filepath.Abs("./fixtures/get-prompt")
	if err != nil {
		t.Fatal(err)
	}

	c, err := scrapligocli.NewCli(
		testHost,
		append(
			[]scrapligooptions.Option{
				scrapligooptions.WithUsername("admin"),
				scrapligooptions.WithPassword("admin"),
				scrapligooptions.WithLookupKeyValue("enable", "libscrapli"),
				scrapligooptions.WithDefinitionFileOrName(scrapligocli.AristaEos),
				scrapligooptions.WithTransportTest(),
				scrapligooptions.WithTestTransportF(f),
				scrapligooptions.WithReadSize(1),
			},
			opts...,
		)...,
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()

	_, err = c.Open(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer closeCancel()

		_, _ = c.Close(closeCtx, scrapligocli.WithForceClose())
	})

	_, err = c.GetPrompt(ctx)
	if err != nil {
		t.Fatal(err)
	}

	return c
}

func TestSendInputErrors(t *testing.T) {
	cases := map[string]struct {
		options  []scrapligooptions.Option
		ctx      func(t *testing.T) context.Context
		expected error
	}{
		"operation-timeout": {
			options: []scrapligooptions.Option{
				scrapligooptions.WithOperationTimeout(time.Second),
			},
			expected: scrapligoerrors.ErrTimeout,
		},
		"context-deadline": {
			ctx: func(t *testing.T) context.Context {
				t.Helper()

				ctx, cancel := context.WithTimeout(t.Context(), time.Second)
				t.Cleanup(cancel)

				return ctx
			},
			expected: context.DeadlineExceeded,
		},
		"context-cancelled": {
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
			d := getUnresponsiveCli(t, c.options...)

			ctx := t.Context()
			if c.ctx != nil {
				ctx = c.ctx(t)
			}

			_, err := d.SendInput(ctx, "show version")
			if !errors.Is(err, c.expected) {
				t.Fatalf("expected errors.Is(err, %v), got %v", c.expected, err)
			}
		})
	}
}
