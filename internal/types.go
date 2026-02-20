package internal

import "encoding/xml"

// NetconfServerHello is an internal struct for unmarshaling server hellos when dealing with
// user provided client capabilities.
type NetconfServerHello struct {
	XMLName      xml.Name                  `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 hello"`
	Capabilities NetconfServerCapabilities `xml:"capabilities"`
}

// NetconfServerCapabilities is an internal struct for unmarshaling server hellos when dealing with
// user provided client capabilities.
type NetconfServerCapabilities struct {
	Capability []string `xml:"capability"`
}
