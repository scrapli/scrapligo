package base_test

import (
	"os"
	"testing"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/google/go-cmp/cmp"
)

func TestTextFsmParse(t *testing.T) {
	r := base.NewResponse("localhost", 22, "show version", []string{})

	f, err := os.ReadFile(
		"../../test_data/driver/network/expected/cisco_iosxe_show_version_textfsm",
	)
	if err != nil {
		t.Fatalf("failed opening channel output file")
	}

	r.Record(f, string(f))

	textfsmOutput, parseErr := r.TextFsmParse(
		"../../test_data/driver/base/response/cisco_ios_show_version.textfsm",
	)
	if parseErr != nil {
		t.Fatalf("failed opening textfsm template file")
	}

	expected := []map[string]interface{}{
		{
			"CONFIG_REGISTER": string("0x2102"),
			"HARDWARE":        []string{"CSR1000V"},
			"HOSTNAME":        string("csr1000v"),
			"MAC":             []string{},
			"RELOAD_REASON":   string("reload"),
			"ROMMON":          string("IOS-XE"),
			"RUNNING_IMAGE":   string("packages.conf"),
			"SERIAL":          []string{"9MVVU09YZFH"},
			"UPTIME":          string("1 hour, 31 minutes"),
			"VERSION":         string("16.12.3"),
		},
	}

	if diff := cmp.Diff(textfsmOutput, expected); diff != "" {
		t.Errorf("actual result and expected result do not match (-want +got):\n%s", diff)
	}
}
