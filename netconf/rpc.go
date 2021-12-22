package netconf

// RPC sends a "bare" rpc to the device.
func (d *Driver) RPC(o ...Option) (*Response, error) {
	finalOpts := d.ParseNetconfOptions(o)
	netconfMessage, err := d.BuildRPCElem(finalOpts.Filter)

	if err != nil {
		return nil, err
	}

	return d.finalizeAndSendMessage(netconfMessage)
}
