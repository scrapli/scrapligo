package netconf

// Commit issue commit rpc to device.
func (d *Driver) Commit() (*Response, error) {
	netconfMessage := d.BuildCommitElem()

	return d.finalizeAndSendMessage(netconfMessage)
}
