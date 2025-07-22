package util

import (
	"os"
	"strconv"
)

// GetEnvStrOrDefault returns the value of the environment variable k as a string *or* the default d
// if casting fails or the environment variable is not set.
func GetEnvStrOrDefault(k, d string) string {
	v, ok := os.LookupEnv(k)
	if ok && v != "" {
		return v
	}

	return d
}

// GetEnvIntOrDefault returns the value of the environment variable k as an int *or* the default d
// if casting fails or the environment variable is not set.
func GetEnvIntOrDefault(k string, d int) int {
	v, ok := os.LookupEnv(k)
	if ok {
		ev, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return d
		}

		return int(ev)
	}

	return d
}
