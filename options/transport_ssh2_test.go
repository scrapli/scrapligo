package options_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	scrapligocli "github.com/kentik/scrapligo/v2/cli"
	scrapligooptions "github.com/kentik/scrapligo/v2/options"
	scrapligotesthelper "github.com/kentik/scrapligo/v2/testhelper"
)

func TestTransportSSH2Options(t *testing.T) {
	d, err := scrapligocli.NewCli(
		"1.2.3.4",
		scrapligooptions.WithTransportSSH2(),
		scrapligooptions.WithSSH2KnownHostsPath("ssh2/known/hosts/path"),
		scrapligooptions.WithSSH2LibSSH2Trace(),
		scrapligooptions.WithSSH2ProxyJumpHost("thisjumpyhost"),
		scrapligooptions.WithSSH2ProxyJumpPort(1234),
		scrapligooptions.WithSSH2ProxyJumpUsername("jumpuser"),
		scrapligooptions.WithSSH2ProxyJumpPassword("jumppass"),
		scrapligooptions.WithSSH2ProxyJumpPrivateKeyPath("jumpprivatekey"),
		scrapligooptions.WithSSH2ProxyJumpPrivateKeyPassphrase("jumpprivatekeypassphrase"),
		scrapligooptions.WithSSH2ProxyJumpLibssh2Trace(),
	)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := d.GetOptions()
	if err != nil {
		t.Fatal(err)
	}

	testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", "transport_ssh2.json"))
	if err != nil {
		t.Fatal(err)
	}

	if *scrapligotesthelper.Update {
		scrapligotesthelper.WriteFile(
			t,
			testGoldenPath,
			[]byte(actual),
		)

		return
	}

	testGoldenContent := string(scrapligotesthelper.ReadFile(t, testGoldenPath))

	if !strings.EqualFold(actual, testGoldenContent) {
		scrapligotesthelper.FailOutput(t, actual, testGoldenContent)
	}
}
