package cfg

func getConfigCommand(configCommandMap map[string]string, source string) (string, error) {
	cmd, ok := configCommandMap[source]

	if !ok {
		return "", ErrInvalidConfigTarget
	}

	return cmd, nil
}
