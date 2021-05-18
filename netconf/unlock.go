package netconf

// Unlock issue unlock rpc to device.
func (d *Driver) Unlock(target string) (*Response, error) {
	netconfMessage := d.BuildUnlockElem(target)

	return d.finalizeAndSendMessage(netconfMessage)
}
