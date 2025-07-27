package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
)

func TestReadWithCallbacks(t *testing.T) {
	testName := "read-with-callbacks"

	testFixturePath, err := filepath.Abs(fmt.Sprintf("./fixtures/%s", testName))
	if err != nil {
		t.Fatal(err)
	}

	testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", testName))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	c := getCli(t, testFixturePath)

	_, err = c.Open(ctx)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_, _ = c.Close(ctx)
	}()

	r, err := c.ReadWithCallbacks(
		ctx,
		"show version",
		scrapligocli.NewReadCallback(
			"cb1",
			func(_ context.Context, c *scrapligocli.Cli) error {
				return c.WriteAndReturn("show version")
			},
			scrapligocli.WithContains("eos1#"),
			scrapligocli.WithOnce(),
		),
		scrapligocli.NewReadCallback(
			"cb2",
			func(_ context.Context, _ *scrapligocli.Cli) error {
				return nil
			},
			scrapligocli.WithContains("eos1#"),
			scrapligocli.WithOnce(),
			scrapligocli.WithCompletes(),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertResult(t, r, testGoldenPath)
}
