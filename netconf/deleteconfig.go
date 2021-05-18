package netconf

// DeleteConfig issue delete rpc to device.
func (d *Driver) DeleteConfig(target string) (*Response, error) {
	netconfMessage := d.BuildDeleteConfigElem(target)

	return d.finalizeAndSendMessage(netconfMessage)
}
