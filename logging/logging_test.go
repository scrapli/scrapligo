package logging_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	scrapligologging "github.com/scrapli/scrapligo/v2/logging"
)

func TestIntFromLevel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		level    scrapligologging.LogLevel
		expected scrapligologging.LogLevelAsInt
	}{
		{
			name:     "trace",
			level:    scrapligologging.Trace,
			expected: scrapligologging.TraceAsInt,
		},
		{
			name:     "debug",
			level:    scrapligologging.Debug,
			expected: scrapligologging.DebugAsInt,
		},
		{
			name:     "info",
			level:    scrapligologging.Info,
			expected: scrapligologging.InfoAsInt,
		},
		{
			name:     "warn",
			level:    scrapligologging.Warn,
			expected: scrapligologging.WarnAsInt,
		},
		{
			name:     "critical",
			level:    scrapligologging.Critical,
			expected: scrapligologging.CriticalAsInt,
		},
		{
			name:     "fatal",
			level:    scrapligologging.Fatal,
			expected: scrapligologging.FatalAsInt,
		},
		{
			name:     "disabled",
			level:    scrapligologging.Disabled,
			expected: scrapligologging.DisabledAsInt,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := scrapligologging.IntFromLevel(tc.level)
			if actual != tc.expected {
				t.Fatalf("IntFromLevel(%q) = %v, expected %v", tc.level, actual, tc.expected)
			}
		})
	}
}

func TestLevelFromInt(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		level    uint8
		expected scrapligologging.LogLevel
	}{
		{
			name:     "trace",
			level:    uint8(scrapligologging.TraceAsInt),
			expected: scrapligologging.Trace,
		},
		{
			name:     "debug",
			level:    uint8(scrapligologging.DebugAsInt),
			expected: scrapligologging.Debug,
		},
		{
			name:     "info",
			level:    uint8(scrapligologging.InfoAsInt),
			expected: scrapligologging.Info,
		},
		{
			name:     "warn",
			level:    uint8(scrapligologging.WarnAsInt),
			expected: scrapligologging.Warn,
		},
		{
			name:     "critical",
			level:    uint8(scrapligologging.CriticalAsInt),
			expected: scrapligologging.Critical,
		},
		{
			name:     "fatal",
			level:    uint8(scrapligologging.FatalAsInt),
			expected: scrapligologging.Fatal,
		},
		{
			name:     "disabled",
			level:    uint8(scrapligologging.DisabledAsInt),
			expected: scrapligologging.Disabled,
		},
		{
			name:     "unknown",
			level:    255,
			expected: scrapligologging.Disabled,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := scrapligologging.LevelFromInt(tc.level)
			if actual != tc.expected {
				t.Fatalf("LevelFromInt(%d) = %q, expected %q", tc.level, actual, tc.expected)
			}
		})
	}
}

func TestLoggerToAnyLoggerStdLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	anyLogger := scrapligologging.LoggerToAnyLogger(logger, scrapligologging.Trace)
	anyLogger.Critical("boom")

	actual := strings.TrimSpace(buf.String())
	expected := "crit :: boom"

	if actual != expected {
		t.Fatalf("output = %q, expected %q", actual, expected)
	}
}
