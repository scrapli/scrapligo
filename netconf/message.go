package netconf

import (
	"encoding/xml"
	"fmt"
)

// BuildFinalMessage creates the final message to send to the device.
func (c *Channel) BuildFinalMessage(xmlMessage interface{}) ([]byte, error) {
	message, err := xml.Marshal(xmlMessage)
	if err != nil {
		return []byte{}, err
	}

	message = append([]byte(XMLHeader), message...)

	if c.SelectedNetconfVersion == Version11 {
		message = append([]byte(fmt.Sprintf("#%d\n", len(message))), message...)
		message = append(message, []byte("\n##")...)
	} else {
		message = append(message, []byte("\n]]>]]>")...)
	}

	return message, nil
}
