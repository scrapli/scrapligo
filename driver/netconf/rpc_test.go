package netconf_test

import (
	"testing"

	"github.com/scrapli/scrapligo/driver/netconf"

	"github.com/google/go-cmp/cmp"
)

func TestForceSelfClosingTags(t *testing.T) {
	tests := map[string]struct {
		got  []byte
		want []byte
	}{
		"empty_tag_no_attrs": {
			got: []byte(
				`<?xml version="1.0" encoding="UTF-8"?><rpc xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" message-id="101"><get-config><source><running></running></source></get-config></rpc>]]>]]>`, //nolint: lll
			),
			want: []byte(
				`<?xml version="1.0" encoding="UTF-8"?><rpc xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" message-id="101"><get-config><source><running/></source></get-config></rpc>]]>]]>`, //nolint: lll
			),
		},
		"empty_tag_with_attrs": {
			got: []byte(
				`<?xml version="1.0" encoding="UTF-8"?><rpc xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" message-id="101"><get><filter type="subtree"><routing-policy xmlns="http://openconfig.net/yang/routing-policy"></routing-policy></filter></get></rpc>]]>]]>`, //nolint: lll
			),
			want: []byte(
				`<?xml version="1.0" encoding="UTF-8"?><rpc xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" message-id="101"><get><filter type="subtree"><routing-policy xmlns="http://openconfig.net/yang/routing-policy"/></filter></get></rpc>]]>]]>`, //nolint: lll
			),
		},
		"empty_tag_with_attrs_and_spaces": {
			got: []byte(
				`<routing-policy xmlns="http://openconfig.net/yang/routing-policy">    </routing-policy>`, //nolint: lll
			),
			want: []byte(`<routing-policy xmlns="http://openconfig.net/yang/routing-policy"/>`),
		},
		"empty_tag_no_attrs_and_spaces": {
			got:  []byte(`<running>  </running>`),
			want: []byte(`<running/>`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := netconf.ForceSelfClosingTags(tt.got)
			if !cmp.Equal(got, tt.want) {
				t.Fatalf(
					"%s: actual and expected values do not match\nactual: %s\nexpected:%s",
					name,
					got,
					tt.want,
				)
			}
		})
	}
}
