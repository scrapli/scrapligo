package base

// ParseSendOptions convenience function to parse and set defaults for `SendOption`s.
func (d *Driver) ParseSendOptions(
	o []SendOption,
) *SendOptions {
	finalOpts := &SendOptions{
		StripPrompt:        DefaultSendOptionsStripPrompt,
		FailedWhenContains: d.FailedWhenContains,
		StopOnFailed:       DefaultSendOptionsStopOnFailed,
		TimeoutOps:         DefaultSendOptionsTimeoutOps,
		Eager:              DefaultSendOptionsEager,
	}

	if len(o) > 0 && o[0] != nil {
		for _, option := range o {
			option(finalOpts)
		}
	}

	return finalOpts
}
