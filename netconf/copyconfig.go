package netconf

// CopyConfig issues copy-config rpc to device.
func (d *Driver) CopyConfig(source, target string) (*Response, error) {
	netconfMessage := d.BuildCopyConfigElem(source, target)

	return d.finalizeAndSendMessage(netconfMessage)
}
