package transport_test

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/logging"
	"github.com/scrapli/scrapligo/transport"
)

func TestSystemTransportNonBlocking(t *testing.T) {
	sendFinished := make(chan struct{})
	netconfHandler := func(s ssh.Session) {
		t.Logf("Handle request")
		time.Sleep(100 * time.Millisecond)
		// io.WriteString(s, "Banner\n\n")
		// send more that 8_192 without a linefeed
		io.WriteString(s, "a")
		// io.WriteString(s, `<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">\n <capabilities>\n  <capability>urn:ietf:params:netconf:base:1.1</capability>\n  <capability>urn:ietf:params:netconf:capability:candidate:1.0</capability>\n  <capability>urn:ietf:params:netconf:capability:rollback-on-error:1.0</capability>\n  <capability>urn:ietf:params:netconf:capability:validate:1.1</capability>\n  <capability>urn:ietf:params:netconf:capability:confirmed-commit:1.1</capability>\n  <capability>urn:ietf:params:netconf:capability:notification:1.0</capability>\n  <capability>urn:ietf:params:netconf:capability:interleave:1.0</capability>\n  <capability>http://cisco.com/ns/yang/Cisco-IOS-XR-infra-systemmib-cfg?module=Cisco-IOS-XR-infra-systemmib-cfg&amp;revision=2015-11-09</capability>\n  <capability>http://cisco.com/ns/yang/Cisco-IOS-XR-ipv4-autorp-datatypes?module=Cisco-IOS-XR-ipv4-autorp-datatypes&amp;revision=2015-11-09</capability>\n  <capability>http://cisco.com/ns/yang/Cisco-IOS-XR-perf-meas-cfg?module=Cisco-IOS-XR-perf-mea`)
		close(sendFinished)
		t.Logf("send finished")
		t.Logf("Request finished")
	}

	sshPort := 2222

	device := dummyDevice(t, sshPort, nil, netconfHandler)
	defer func() {
		if err := device.Close(); err != nil {
			t.Error(err)
		}
	}()

	sshArgs, err := transport.NewSSHArgs(
		options.WithAuthNoStrictKey(),
		options.WithSSHKnownHostsFile("/dev/null"),
	)
	if err != nil {
		t.Fatal(err)
	}

	sshArgs.NetconfConnection = true

	tp, err := transport.NewSystemTransport(sshArgs)
	if err != nil {
		t.Fatal(err)
	}

	openArgs, err := transport.NewArgs(
		&logging.Instance{},
		"localhost",
		options.WithPort(sshPort),
		options.WithAuthUsername("whatever"),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = tp.Open(openArgs)
	if err != nil {
		t.Fatal(err)
	}

	doneChan := make(chan struct{})
	go func() {
		defer t.Log("read finished")
		defer close(doneChan)
		for {
			t.Log("starting to read")
			// from the channel code: defaultReadSize = 8_192
			b, err := tp.Read(1)
			t.Logf("read %d bytes: %s", len(b), b)
			if err != nil {
				if err == io.EOF {
					return
				}
				t.Logf("failed to read : %s", err)
				return
			}
		}
	}()
	<-sendFinished
	time.Sleep(5 * time.Second)

	t.Log("closing transport")
	err = tp.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("waiting for read to be done")
	<-doneChan
}

func TestSystemTransportBlocking(t *testing.T) {
	sendFinished := make(chan struct{})
	netconfHandler := func(s ssh.Session) {
		t.Logf("Handle request")
		time.Sleep(100 * time.Millisecond)
		// io.WriteString(s, "Banner\n\n")
		// send more that 8_192 without a linefeed
		io.WriteString(s, strings.Repeat("z", 8_192))
		io.WriteString(s, strings.Repeat("a", 10))
		// io.WriteString(s, `<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">\n <capabilities>\n  <capability>urn:ietf:params:netconf:base:1.1</capability>\n  <capability>urn:ietf:params:netconf:capability:candidate:1.0</capability>\n  <capability>urn:ietf:params:netconf:capability:rollback-on-error:1.0</capability>\n  <capability>urn:ietf:params:netconf:capability:validate:1.1</capability>\n  <capability>urn:ietf:params:netconf:capability:confirmed-commit:1.1</capability>\n  <capability>urn:ietf:params:netconf:capability:notification:1.0</capability>\n  <capability>urn:ietf:params:netconf:capability:interleave:1.0</capability>\n  <capability>http://cisco.com/ns/yang/Cisco-IOS-XR-infra-systemmib-cfg?module=Cisco-IOS-XR-infra-systemmib-cfg&amp;revision=2015-11-09</capability>\n  <capability>http://cisco.com/ns/yang/Cisco-IOS-XR-ipv4-autorp-datatypes?module=Cisco-IOS-XR-ipv4-autorp-datatypes&amp;revision=2015-11-09</capability>\n  <capability>http://cisco.com/ns/yang/Cisco-IOS-XR-perf-meas-cfg?module=Cisco-IOS-XR-perf-mea`)
		close(sendFinished)
		t.Logf("send finished")
		time.Sleep(60 * time.Second)
		t.Logf("Request finished")
	}

	sshPort := 2222

	device := dummyDevice(t, sshPort, nil, netconfHandler)
	defer func() {
		if err := device.Close(); err != nil {
			t.Error(err)
		}
	}()

	sshArgs, err := transport.NewSSHArgs(
		options.WithAuthNoStrictKey(),
		options.WithSSHKnownHostsFile("/dev/null"),
	)
	if err != nil {
		t.Fatal(err)
	}

	sshArgs.NetconfConnection = true

	tp, err := transport.NewSystemTransport(sshArgs)
	if err != nil {
		t.Fatal(err)
	}

	openArgs, err := transport.NewArgs(
		&logging.Instance{},
		"localhost",
		options.WithPort(sshPort),
		options.WithAuthUsername("whatever"),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = tp.Open(openArgs)
	if err != nil {
		t.Fatal(err)
	}

	doneChan := make(chan struct{})
	go func() {
		defer t.Log("read finished")
		defer close(doneChan)
		for {
			t.Log("starting to read")
			// from the channel code: defaultReadSize = 8_192
			b, err := tp.Read(81)
			t.Logf("read %d bytes: %s", len(b), b)
			if err != nil {
				if err == io.EOF {
					return
				}
				t.Logf("failed to read : %s", err)
				return
			}
		}
	}()
	<-sendFinished
	time.Sleep(5 * time.Second)

	t.Log("closing transport")
	err = tp.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("waiting for read to be done")
	<-doneChan
}
