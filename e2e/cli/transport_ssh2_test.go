package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestSSH2TransportProxyJump(t *testing.T) {
	parentName := "ssh2-transport-proxy-jump"

	cases := map[string]struct {
		description string
		platform    string
		input       string
	}{
		"eos": {
			description: "simple",
			platform:    scrapligocli.AristaEos.String(),
			input:       "show version | i Kern",
		},
		"srl": {
			description: "simple",
			platform:    scrapligocli.NokiaSrlinux.String(),
			input:       "show version | grep OS",
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
				scrapligooptions.WithDefintionFileOrName(caseData.platform),
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
					scrapligooptions.WithSSH2ProxyJumpUsername("admin"),
					scrapligooptions.WithSSH2ProxyJumpPassword("NokiaSrl1!"),
				)
			}

			c, err := scrapligocli.NewCli(
				host,
				opts...,
			)
			if err != nil {
				t.Fatal(err)
			}

			_, err = c.Open(ctx)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_, _ = c.Close(ctx)
			}()

			r, err := c.SendInput(ctx, caseData.input)
			if err != nil {
				t.Fatal(err)
			}

			assertResult(t, r, testGoldenPath)
		})
	}
}
