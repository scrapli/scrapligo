package cfg

import (
	"regexp"
	"strconv"

	"github.com/scrapli/scrapligo/driver/base"
)

const (
	execPrivLevel            = "privilege_exec"
	configExclusivePrivLevel = "configuration_exclusive"
)

func getConfigCommand(configCommandMap map[string]string, source string) (string, error) {
	cmd, ok := configCommandMap[source]

	if !ok {
		return "", ErrInvalidConfigTarget
	}

	return cmd, nil
}

func parseSpaceAvail(
	bytesFreePattern *regexp.Regexp,
	filesystemSizeResult *base.Response,
) (int, error) {
	var err error

	bytesAvailMatch := bytesFreePattern.FindStringSubmatch(filesystemSizeResult.Result)

	bytesAvail := -1

	for i, name := range bytesFreePattern.SubexpNames() {
		if i != 0 && name == "bytes_available" {
			bytesAvail, err = strconv.Atoi(bytesAvailMatch[i])
			if err != nil {
				return -1, err
			}
		}
	}

	return bytesAvail, nil
}

func isSpaceSufficient(
	filesystemBytesAvail int,
	filesystemSpaceAvailBufferPerc float32,
	config string,
) bool {
	return float32(
		filesystemBytesAvail,
	) >= float32(
		len(config),
	)/(filesystemSpaceAvailBufferPerc/100.0)+float32(
		len(config),
	)
}
