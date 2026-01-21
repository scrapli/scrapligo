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

func TestBinTransportProxyJump(t *testing.T) {
	parentName := "bin-transport-proxy-jump"

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

			sshConfigFilePath := "./fixtures/ssh_config"
			if runtime.GOOS == scrapligoconstants.Darwin {
				sshConfigFilePath += "_darwin"
			} else {
				sshConfigFilePath += "_linux"
			}

			opts := []scrapligooptions.Option{
				scrapligooptions.WithDefinitionFileOrName(caseData.platform),
				scrapligooptions.WithUsername("admin"),
				scrapligooptions.WithTransportBin(),
				scrapligooptions.WithBinTransportSSHConfigFile(sshConfigFilePath),
			}

			var host string

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
