package netconf

import (
	"encoding/xml"
	"fmt"
)

type message struct {
	XMLName   xml.Name    `xml:"rpc"`
	Namespace string      `xml:"xmlns,attr"`
	MessageID int         `xml:"message-id,attr"`
	Payload   interface{} `xml:",innerxml"`
}

type serializedInput struct {
	rawXML    []byte
	framedXML []byte
}

func (m *message) serialize(v string, forceSelfClosingTags bool) (*serializedInput, error) {
	serialized := &serializedInput{}

	msg, err := xml.Marshal(m)
	if err != nil {
		return nil, err
	}

	msg = append([]byte(xmlHeader), msg...)

	if forceSelfClosingTags {
		msg = ForceSelfClosingTags(msg)
	}

	// copy the raw xml (without the netconf framing) before setting up framing bits
	serialized.rawXML = make([]byte, len(msg))
	copy(serialized.rawXML, msg)

	switch v {
	case V1Dot0:
		msg = append(msg, []byte(v1Dot0Delim)...)
	case V1Dot1:
		msg = append([]byte(fmt.Sprintf("#%d\n", len(msg))), msg...)
		msg = append(msg, []byte("\n##")...)
	}

	serialized.framedXML = msg

	return serialized, nil
}
