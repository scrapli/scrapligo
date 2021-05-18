package netconf

// Get execute get rpc with optional filters.
func (d *Driver) Get(o ...Option) (*Response, error) {
	finalOpts := d.ParseNetconfOptions(o)
	netconfMessage, err := d.BuildGetElem(finalOpts.Filter, finalOpts.FilterType)

	if err != nil {
		return nil, err
	}

	return d.finalizeAndSendMessage(netconfMessage)
}
