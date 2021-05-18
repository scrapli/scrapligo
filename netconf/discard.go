package netconf

// Discard issue discard rpc to device.
func (d *Driver) Discard() (*Response, error) {
	netconfMessage := d.BuildDiscardElem()

	return d.finalizeAndSendMessage(netconfMessage)
}
