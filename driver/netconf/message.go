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

func (m *message) serialize(v string, forceSelfClosingTags bool) ([]byte, error) {
	msg, err := xml.Marshal(m)
	if err != nil {
		return nil, err
	}

	msg = append([]byte(xmlHeader), msg...)

	if forceSelfClosingTags {
		msg = ForceSelfClosingTags(msg)
	}

	switch v {
	case V1Dot0:
		msg = append(msg, []byte(v1Dot0Delim)...)
	case V1Dot1:
		msg = append([]byte(fmt.Sprintf("#%d\n", len(msg))), msg...)
		msg = append(msg, []byte("\n##")...)
	}

	return msg, nil
}
