package cli_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	scrapligocli "github.com/scrapli/scrapligo/cli"
	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
)

func TestOpenWithKey(t *testing.T) {
	parentName := "open-with-key"

	cases := map[string]struct {
		description string
		options     []scrapligooptions.Option
	}{
		"bin": {
			description: "simple-bin",
			options: []scrapligooptions.Option{
				scrapligooptions.WithTransportBin(),
				scrapligooptions.WithUsername("admin-sshkey"),
				scrapligooptions.WithPrivateKeyPath("./fixtures/libscrapli_test_ssh_key"),
				scrapligooptions.WithLookupKeyValue("enable", "libscrapli"),
			},
		},
		"ssh2": {
			description: "simple-ssh2",
			options: []scrapligooptions.Option{
				scrapligooptions.WithTransportSSH2(),
				scrapligooptions.WithUsername("admin-sshkey"),
				scrapligooptions.WithPrivateKeyPath("./fixtures/libscrapli_test_ssh_key"),
				scrapligooptions.WithLookupKeyValue("enable", "libscrapli"),
			},
		},
		"bin-passhprase": {
			description: "bin-with-passhrase",
			options: []scrapligooptions.Option{
				scrapligooptions.WithTransportBin(),
				scrapligooptions.WithUsername("admin-sshkey-passphrase"),
				scrapligooptions.WithPrivateKeyPath(
					"./fixtures/libscrapli_test_ssh_key_passphrase",
				),
				scrapligooptions.WithPrivateKeyPassphrase("libscrapli"),
				scrapligooptions.WithLookupKeyValue("enable", "libscrapli"),
			},
		},
		"ssh2-passhrase": {
			description: "ssh2-with-passhrase",
			options: []scrapligooptions.Option{
				scrapligooptions.WithTransportSSH2(),
				scrapligooptions.WithUsername("admin-sshkey-passphrase"),
				scrapligooptions.WithPrivateKeyPath(
					"./fixtures/libscrapli_test_ssh_key_passphrase",
				),
				scrapligooptions.WithPrivateKeyPassphrase("libscrapli"),
				scrapligooptions.WithLookupKeyValue("enable", "libscrapli"),
			},
		},
	}

	for caseName, caseData := range cases {
		testName := fmt.Sprintf("%s-%s", parentName, caseName)

		t.Run(testName, func(t *testing.T) {
			t.Logf("%s: starting", testName)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			var host string

			if !scrapligotesthelper.EosAvailable() {
				t.Skip("skipping case, arista eos unavailable...")
			}

			var port uint16

			if runtime.GOOS == scrapligoconstants.Darwin {
				host = localhost
				port = 22022
			} else {
				host = "172.20.20.17"
			}

			caseData.options = append(
				caseData.options,
				scrapligooptions.WithPort(port),
			)

			c, err := scrapligocli.NewCli(
				scrapligocli.AristaEos,
				host,
				caseData.options...,
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

			// no need to assert result, just make sure getprompt works
			_, err = c.GetPrompt(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
