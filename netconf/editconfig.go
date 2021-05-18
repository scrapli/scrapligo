package netconf

// EditConfig edit device configuration.
func (d *Driver) EditConfig(target, config string) (*Response, error) {
	netconfMessage := d.BuildEditConfigElem(config, target)

	return d.finalizeAndSendMessage(netconfMessage)
}
