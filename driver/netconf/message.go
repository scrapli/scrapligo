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
	message, err := xml.Marshal(m)
	if err != nil {
		return nil, err
	}

	message = append([]byte(xmlHeader), message...)

	if forceSelfClosingTags {
		message = ForceSelfClosingTags(message)
	}

	switch v {
	case V1Dot0:
		message = append(message, []byte(v1Dot0Delim)...)
	case V1Dot1:
		message = append([]byte(fmt.Sprintf("#%d\n", len(message))), message...)
		message = append(message, []byte("\n##")...)
	}

	return message, nil
}
