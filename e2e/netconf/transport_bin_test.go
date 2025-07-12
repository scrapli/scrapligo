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

const localhost = "localhost"

func TestBinTransportProxyJumpNetconf(t *testing.T) {
	parentName := "bin-transport-proxy-jump-netconf"

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
				scrapligooptions.WithUsername("admin"),
			}

			if caseData.platform == scrapligocli.AristaEos.String() {
				if !scrapligotesthelper.EosAvailable() {
					t.Skip("skipping case, arista eos unavailable...")
				}

				opts = append(
					opts,
					scrapligooptions.WithPassword("admin"),
					scrapligooptions.WithLookupKeyValue("enable", "libscrapli"),
				)

				var port uint16

				if runtime.GOOS == scrapligoconstants.Darwin {
					host = localhost
					port = 22830
				} else {
					host = "172.20.20.17"
				}

				opts = append(
					opts,
					scrapligooptions.WithPort(port),
				)
			} else {
				opts = append(
					opts,
					scrapligooptions.WithPassword("NokiaSrl1!"),
				)

				if runtime.GOOS == scrapligoconstants.Darwin {
					host = localhost

					opts = append(
						opts,
						scrapligooptions.WithPort(21830),
					)
				} else {
					host = "172.20.20.16"
				}
			}

			sshConfigFilePath := "./fixtures/ssh_config"
			if runtime.GOOS == scrapligoconstants.Darwin {
				sshConfigFilePath += "_darwin"
			} else {
				sshConfigFilePath += "_linux"
			}

			opts = append(
				opts,
				scrapligooptions.WithTransportBin(),
				scrapligooptions.WithBinTransportSSHConfigFile(sshConfigFilePath),
			)

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
