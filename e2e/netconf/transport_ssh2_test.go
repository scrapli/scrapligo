package netconf_test

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligonetconf "github.com/scrapli/scrapligo/netconf"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestSSH2TransportProxyJumpNetconf(t *testing.T) {
	parentName := "ssh2-transport-proxy-jump-netconf"

	cases := map[string]struct {
		description string
		platform    string
		filter      string
	}{
		"eos": {
			description: "simple",
			platform:    scrapligocli.AristaEos.String(),
			filter:      "<system><config><hostname></hostname></config></system>",
		},
		"srl": {
			description: "simple",
			platform:    scrapligocli.NokiaSrlinux.String(),
			filter:      "<system xmlns=\"urn:nokia.com:srlinux:general:system\"><ssh-server xmlns=\"urn:nokia.com:srlinux:linux:ssh\"><name>mgmt</name></ssh-server></system>",
		},
	}

	for caseName, caseData := range cases {
		testName := fmt.Sprintf("%s-%s", parentName, caseName)

		t.Run(testName, func(t *testing.T) {
			t.Logf("%s: starting", testName)

			testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", testName))
			if err != nil {
				t.Fatal(err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			var host string

			opts := []scrapligooptions.Option{
				scrapligooptions.WithUsername("scrapli-pw"),
				scrapligooptions.WithPassword("scrapli-123-pw"),
				scrapligooptions.WithTransportSSH2(),
			}

			if runtime.GOOS == scrapligoconstants.Darwin {
				host = localhost

				opts = append(
					opts,
					scrapligooptions.WithPort(24022),
				)
			} else {
				host = "172.20.20.19"
			}

			if caseData.platform == scrapligocli.AristaEos.String() {
				if !scrapligotesthelper.EosAvailable() {
					t.Skip("skipping case, arista eos unavailable...")
				}

				opts = append(
					opts,
					scrapligooptions.WithSSH2ProxyJumpHost("172.20.20.17"),
					scrapligooptions.WithSSH2ProxyJumpUsername("admin"),
					scrapligooptions.WithSSH2ProxyJumpPassword("admin"),
					scrapligooptions.WithLookupKeyValue("enable", "libscrapli"),
				)
			} else {
				opts = append(
					opts,
					scrapligooptions.WithSSH2ProxyJumpHost("172.20.20.16"),
					scrapligooptions.WithSSH2ProxyJumpPort(830),
					scrapligooptions.WithSSH2ProxyJumpUsername("admin"),
					scrapligooptions.WithSSH2ProxyJumpPassword("NokiaSrl1!"),
				)
			}

			n, err := scrapligonetconf.NewNetconf(
				host,
				opts...,
			)
			if err != nil {
				t.Fatal(err)
			}

			_, err = n.Open(ctx)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_, _ = n.Close(ctx)
			}()

			r, err := n.Get(ctx, scrapligonetconf.WithFilter(caseData.filter))
			if err != nil {
				t.Fatal(err)
			}

			assertResult(t, r, testGoldenPath)
		})
	}
}
