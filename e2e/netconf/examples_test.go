package netconf_test

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestNetconfExamples(t *testing.T) {
	parentName := "examples-netconf"

	examples, err := filepath.Glob("../../examples/netconf/*")
	if err != nil {
		t.Fatalf("failed globbing cli examples, error: %v", err)
	}

	for _, example := range examples {
		testName := fmt.Sprintf("%s-%s", parentName, filepath.Base(example))

		t.Run(testName, func(t *testing.T) {
			t.Logf("%s: starting", testName)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			c := exec.CommandContext( //nolint: gosec
				ctx,
				"go",
				"run",
				filepath.Join("examples", "cli", filepath.Base(example), "main.go"),
			)
			c.Dir = "../../"

			err := c.Run()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
