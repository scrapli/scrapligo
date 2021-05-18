package netconf

// Validate issue validate rpc to device.
func (d *Driver) Validate(target string) (*Response, error) {
	netconfMessage := d.BuildValidateElem(target)

	return d.finalizeAndSendMessage(netconfMessage)
}
