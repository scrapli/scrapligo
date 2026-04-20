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

func TestAuthOptions(t *testing.T) {
	d, err := scrapligocli.NewCli(
		"1.2.3.4",
		scrapligooptions.WithUsername("foo"),
		scrapligooptions.WithPassword("bar"),
		scrapligooptions.WithPrivateKeyPath("secretpath"),
		scrapligooptions.WithPrivateKeyPassphrase("secretpassphrase"),
		scrapligooptions.WithLookupKeyValue("baz", "qux"),
		scrapligooptions.WithForceInSessionAuth(),
		scrapligooptions.WithUsernamePattern("usernamepattern"),
		scrapligooptions.WithPasswordPattern("passwordpattern"),
		scrapligooptions.WithPassphrasePattern("passphrasepattern"),
	)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := d.GetOptions()
	if err != nil {
		t.Fatal(err)
	}

	testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", "auth.json"))
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
