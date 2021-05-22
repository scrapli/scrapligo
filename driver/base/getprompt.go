package base

// GetPrompt fetch device prompt.
func (d *Driver) GetPrompt() (string, error) {
	return d.Channel.GetPrompt()
}
