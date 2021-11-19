package channel_test

import (
	"testing"

	"github.com/scrapli/scrapligo/channel"
)

// TestGetAuthPatternsRace ensures that the global singleton-y auth patterns getter is goroutine
// safe.
func TestGetAuthPatternsRace(t *testing.T) {
	for i := 1; i < 5; i++ {
		go channel.GetAuthPatterns()
	}
}
