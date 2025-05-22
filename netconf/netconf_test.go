package netconf_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	scrapligoffi "github.com/scrapli/scrapligo/ffi"
	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

const (
	testHost = "localhost"
)

func TestMain(m *testing.M) {
	scrapligotesthelper.Flags()

	exitCode := m.Run()

	if scrapligoffi.AssertNoLeaks() != nil {
		_, _ = fmt.Fprintln(os.Stderr, "memory leak(s) detected!")

		os.Exit(127)
	}

	_, _ = fmt.Fprintln(os.Stderr, "no memory leak(s) detected!")

	os.Exit(exitCode)
}

func getNetconf(t *testing.T, f string) *scrapligonetconf.Netconf {
	t.Helper()

	opts := []scrapligooptions.Option{
		// note that netconf-admin bypasses enable secret stuff, without this was getting
		// permission denied committing things and such... but wanted to retain the enable
		// secret stuff since its nice to validate default mode gets acquired and stuff
		scrapligooptions.WithUsername("root"),
		scrapligooptions.WithPassword("password"),
		scrapligooptions.WithPort(23830),
	}

	if *scrapligotesthelper.Record {
		opts = append(
			opts,
			scrapligooptions.WithSessionRecorderPath(f),
		)
	} else {
		opts = append(
			opts,
			scrapligooptions.WithTransportTest(),
			scrapligooptions.WithTestTransportF(f),
			scrapligooptions.WithReadSize(1),
			// see libscrapli notes in integration netconf tests
			scrapligooptions.WithOperationMaxSearchDepth(32),
		)
	}

	d, err := scrapligonetconf.NewNetconf(
		testHost,
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func getNetconfSrl(t *testing.T, f string) *scrapligonetconf.Netconf {
	t.Helper()

	opts := []scrapligooptions.Option{
		scrapligooptions.WithUsername("admin"),
		scrapligooptions.WithPassword("NokiaSrl1!"),
		scrapligooptions.WithPort(21830),
	}

	if *scrapligotesthelper.Record {
		opts = append(
			opts,
			scrapligooptions.WithSessionRecorderPath(f),
		)
	} else {
		opts = append(
			opts,
			scrapligooptions.WithTransportTest(),
			scrapligooptions.WithTestTransportF(f),
			scrapligooptions.WithReadSize(1),
			// see libscrapli notes in integration netconf tests
			scrapligooptions.WithOperationMaxSearchDepth(32),
		)
	}

	d, err := scrapligonetconf.NewNetconf(
		testHost,
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func closeNetconf(t *testing.T, n *scrapligonetconf.Netconf) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = n.Close(ctx, scrapligonetconf.WithForceClose())
}

func assertResult(t *testing.T, r *scrapligonetconf.Result, testGoldenPath string) {
	t.Helper()

	if *scrapligotesthelper.Update {
		scrapligotesthelper.WriteFile(
			t,
			testGoldenPath,
			scrapligotesthelper.CleanNetconfOutput(t, r.Result),
		)

		return
	}

	cleanedActual := scrapligotesthelper.CleanNetconfOutput(t, r.Result)

	// we can't just write the cleaned stuff to disk because then chunk sizes will be wrong if we
	// just do the lazy cleanup method we are doing (and cant stop wont stop)
	testGoldenContent := scrapligotesthelper.ReadFile(t, testGoldenPath)
	cleanedGolden := scrapligotesthelper.CleanNetconfOutput(t, string(testGoldenContent))

	if !bytes.Equal(cleanedActual, cleanedGolden) {
		scrapligotesthelper.FailOutput(t, cleanedActual, cleanedGolden)
	}

	scrapligotesthelper.AssertNotDefault(t, r.StartTime)
	scrapligotesthelper.AssertNotDefault(t, r.EndTime)
	scrapligotesthelper.AssertNotDefault(t, r.ElapsedTimeSeconds)
	scrapligotesthelper.AssertNotDefault(t, r.Host)
	scrapligotesthelper.AssertNotDefault(t, r.ResultRaw)
	scrapligotesthelper.AssertEqual(t, false, r.Failed)
}

func TestGetSessionID(t *testing.T) {
	testName := "get-session-id"

	testFixturePath, err := filepath.Abs(fmt.Sprintf("./fixtures/%s", testName))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	n := getNetconf(t, testFixturePath)

	_, err = n.Open(ctx)
	if err != nil {
		t.Fatal(err)
	}

	defer closeNetconf(t, n)

	actual, err := n.GetSessionID()
	if err != nil {
		t.Fatal(err)
	}

	if actual == 0 {
		t.Fatal("expected sesion id to be non zero")
	}
}
