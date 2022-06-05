package netconf

import (
	"encoding/xml"

	"github.com/scrapli/scrapligo/response"
)

type commit struct {
	XMLName xml.Name `xml:"commit"`
}

func (d *Driver) buildCommitElem() *message {
	commitElem := &commit{
		XMLName: xml.Name{},
	}

	netconfInput := d.buildPayload(commitElem)

	return netconfInput
}

// Commit executes a commit rpc against the NETCONF server.
func (d *Driver) Commit() (*response.NetconfResponse, error) {
	op, err := NewOperation()
	if err != nil {
		return nil, err
	}

	return d.sendRPC(d.buildCommitElem(), op)
}

type discard struct {
	XMLName xml.Name `xml:"discard-changes"`
}

func (d *Driver) buildDiscardElem() *message {
	discardElem := &discard{
		XMLName: xml.Name{},
	}

	netconfInput := d.buildPayload(discardElem)

	return netconfInput
}

// Discard executes a discard rpc against the NETCONF server.
func (d *Driver) Discard() (*response.NetconfResponse, error) {
	op, err := NewOperation()
	if err != nil {
		return nil, err
	}

	return d.sendRPC(d.buildDiscardElem(), op)
}
