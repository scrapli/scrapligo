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

const localhost = "localhost"

func TestBinTransportProxyJump(t *testing.T) {
	parentName := "bin-transport-proxy-jump"

	cases := map[string]struct {
		description string
		platform    string
		input       string
		options     []scrapligocli.Option
	}{
		"eos": {
			description: "simple",
			platform:    scrapligocli.AristaEos.String(),
			input:       "show version | i Kern",
			options:     []scrapligocli.Option{},
		},
		"srl": {
			description: "simple",
			platform:    scrapligocli.NokiaSrl.String(),
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
					port = 22022
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
						scrapligooptions.WithPort(21022),
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

			c, err := scrapligocli.NewCli(
				caseData.platform,
				host,
				opts...,
			)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_, _ = c.Close(ctx)
			}()

			_, err = c.Open(ctx)
			if err != nil {
				t.Fatal(err)
			}

			r, err := c.SendInput(ctx, caseData.input, caseData.options...)
			if err != nil {
				t.Fatal(err)
			}

			assertResult(t, r, testGoldenPath)
		})
	}
}
