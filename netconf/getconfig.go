package netconf

// GetConfig execute get-config rpc with optional filters.
func (d *Driver) GetConfig(source string, o ...Option) (*Response, error) {
	finalOpts := d.ParseNetconfOptions(o)

	netconfMessage, err := d.BuildGetConfigElem(
		source,
		finalOpts.Filter,
		finalOpts.FilterType,
		finalOpts.DefaultType,
	)

	if err != nil {
		return nil, err
	}

	return d.finalizeAndSendMessage(netconfMessage)
}
