package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

func TestGetPrompt(t *testing.T) {
	testName := "get-prompt"

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

	r, err := c.GetPrompt(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assertResult(t, r, testGoldenPath)
}
