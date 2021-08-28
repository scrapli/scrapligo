package netconf

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"regexp"
)

// BuildFinalMessage creates the final message to send to the device.
func (c *Channel) BuildFinalMessage(xmlMessage interface{}) ([]byte, error) {
	message, err := xml.Marshal(xmlMessage)
	if err != nil {
		return []byte{}, err
	}

	message = append([]byte(XMLHeader), message...)

	if c.ForceSelfClosingTag {
		// this grossness is 100% only for junos who seem to have lost their ever loving mind...
		// `<source><running></running></source>` will cause junos (at least 17.x) to return an
		// error whilst `<source><running/></source>` will not. functionally 100% identical but
		// sure juniper do you or whatever...
		p := regexp.MustCompile(`(?:<source><(\w+)>.*</source>)`)
		o := p.FindAllSubmatch(message, -1)

		if len(o) > 0 {
			newSource := []byte(fmt.Sprintf("<source><%s/></source>", o[0][1]))
			message = bytes.Replace(message, o[0][0], newSource, 1)
		}
	}

	if c.SelectedNetconfVersion == Version11 {
		message = append([]byte(fmt.Sprintf("#%d\n", len(message))), message...)
		message = append(message, []byte("\n##")...)
	} else {
		message = append(message, []byte("]]>]]>")...)
	}

	return message, nil
}
