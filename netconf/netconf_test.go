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

	n, err := scrapligonetconf.NewNetconf(
		testHost,
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return n
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

	n, err := scrapligonetconf.NewNetconf(
		testHost,
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	return n
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

	// we dont check failed since for now some things (cancel commit) fail expectedly, but we are
	// more just making sure the rpc was successful and we sent valid stuff etc.
	scrapligotesthelper.AssertNotDefault(t, r.StartTime)
	scrapligotesthelper.AssertNotDefault(t, r.EndTime)
	scrapligotesthelper.AssertNotDefault(t, r.ElapsedTimeSeconds)
	scrapligotesthelper.AssertNotDefault(t, r.Host)
	scrapligotesthelper.AssertNotDefault(t, r.ResultRaw)
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

	defer func() {
		_, _ = n.Close(ctx)
	}()

	actual, err := n.GetSessionID()
	if err != nil {
		t.Fatal(err)
	}

	if actual == 0 {
		t.Fatal("expected sesion id to be non zero")
	}
}

func TestGetNextNotification(t *testing.T) {
	testName := "get-next-notification"

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

	defer func() {
		_, _ = n.Close(ctx)
	}()

	_, err = n.RawRPC(
		ctx,
		`<create-subscription xmlns="urn:ietf:params:xml:ns:netconf:notification:1.0">
			<stream>NETCONF</stream>
			<filter type="subtree">
				<counter-update xmlns="urn:boring:counter"/>
			</filter>
		</create-subscription>`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if *scrapligotesthelper.Record {
		// boring counter updates every 3s; obviously only when not using fixture
		time.Sleep(4 * time.Second)
	}

	notif, err := n.GetNextNotification()
	if err != nil {
		t.Fatal(err)
	}

	scrapligotesthelper.AssertNotDefault(t, notif)
}

func TestGetNextSubscription(t *testing.T) {
	testName := "get-next-subscription"

	testFixturePath, err := filepath.Abs(fmt.Sprintf("./fixtures/%s", testName))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	opts := []scrapligooptions.Option{
		scrapligooptions.WithUsername("SOMEUSER"),
		scrapligooptions.WithPassword("SOMEPASSWORD"),
		scrapligooptions.WithPort(830),
	}

	if *scrapligotesthelper.Record {
		t.Fatal( //nolint: revive
			// if you see this and are actually going to update stuff, user/pass above and
			// host below plz kthxbye
			"are you really sure? this is not using the clab setup, " +
				"make sure you have either cisco sandbox or something else " +
				"handy to re-record this test fixture",
		)

		opts = append(
			opts,
			scrapligooptions.WithSessionRecorderPath(testFixturePath),
		)
	} else {
		opts = append(
			opts,
			scrapligooptions.WithTransportTest(),
			scrapligooptions.WithTestTransportF(testFixturePath),
			scrapligooptions.WithReadSize(1),
			// see libscrapli notes in integration netconf tests
			scrapligooptions.WithOperationMaxSearchDepth(32),
		)
	}

	n, err := scrapligonetconf.NewNetconf(
		"HOST",
		opts...,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = n.Open(ctx)
	if err != nil {
		t.Fatal(err)
	}

	r, err := n.RawRPC(
		ctx,
		`<establish-subscription xmlns="urn:ietf:params:xml:ns:yang:ietf-event-notifications" xmlns:yp="urn:ietf:params:xml:ns:yang:ietf-yang-push">
			<stream>yp:yang-push</stream>
			<yp:xpath-filter>/mdt-oper:mdt-oper-data/mdt-subscriptions</yp:xpath-filter>
			<yp:period>1000</yp:period>
		</establish-subscription>`,
	)
	if err != nil {
		t.Fatal(err)
	}

	subID, err := n.GetSubscriptionID(r.Result)
	if err != nil {
		t.Fatal(err)
	}

	if *scrapligotesthelper.Record {
		// obviously only have to wait w/ a real device
		time.Sleep(10 * time.Second)
	}

	sub, err := n.GetNextSubscription(subID)
	if err != nil {
		t.Fatal(err)
	}

	scrapligotesthelper.AssertNotDefault(t, sub)
}
