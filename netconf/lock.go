package netconf

// Lock issue lock rpc to device.
func (d *Driver) Lock(target string) (*Response, error) {
	netconfMessage := d.BuildLockElem(target)

	return d.finalizeAndSendMessage(netconfMessage)
}
