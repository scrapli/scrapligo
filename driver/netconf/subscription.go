package netconf

import (
	"encoding/xml"
	"strconv"

	"github.com/scrapli/scrapligo/response"
)

type establishSubscription struct {
	XMLName               xml.Name `xml:"establish-subscription"`
	NamespaceNotification string   `xml:"xmlns,attr"`
	NamespaceYANGPush     string   `xml:"xmlns:yp,attr"`
	Stream                string   `xml:"stream"`
	Filter                string   `xml:"yp:xpath-filter"`
	Period                int      `xml:"yp:period"`
}

// EstablishPeriodicSubscription is a BETA method to establish a NETCONF subscription. Seriously,
// don't trust that this won't change, just pretend it doesn't exist for now or something!
func (d *Driver) EstablishPeriodicSubscription(
	xpath string,
	period int,
) (*response.NetconfResponse, error) {
	establishElem := &establishSubscription{
		XMLName:               xml.Name{},
		NamespaceNotification: "urn:ietf:params:xml:ns:yang:ietf-event-notifications",
		NamespaceYANGPush:     "urn:ietf:params:xml:ns:yang:ietf-yang-push",
		Stream:                "yp:yang-push",
		Filter:                xpath,
		Period:                period,
	}

	m := d.buildPayload(establishElem)

	r, err := d.sendRPC(m, &OperationOptions{})
	if err != nil {
		return nil, err
	}

	patterns := getNetconfPatterns()

	match := patterns.subscriptionID.FindSubmatch(r.RawResult)
	subID, _ := strconv.Atoi(string(match[1]))

	d.subscriptions[subID] = make([][]byte, 0)

	r.SubscriptionID = subID

	return r, nil
}
