package cfg

import (
	"regexp"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

type EOSCfg struct {
	conn             *network.Driver
	VersionPattern   *regexp.Regexp
	configCommandMap map[string]string
}

// NewEOSCfg return a cfg instance setup for an Arista EOS device.
func NewEOSCfg(
	conn *network.Driver,
	options ...Option,
) (*Cfg, error) {
	options = append([]Option{WithConfigSources([]string{"running", "startup"})}, options...)

	c, err := newCfg(conn, options...)
	if err != nil {
		return nil, err
	}

	c.Platform = &EOSCfg{
		conn:           conn,
		VersionPattern: regexp.MustCompile(`(?i)\d+\.\d+\.[a-z0-9\-]+(\.\d+[a-z]?)?`),
	}

	err = setPlatformOptions(c.Platform, options...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// GetVersion get the version from the device.
func (p *EOSCfg) GetVersion() *Response {
	r := NewResponse(p.conn.Host, nil)

	versionResult, err := p.conn.SendCommand("show version | i Software image version")
	if err != nil {
		return r
	}

	r.Record([]*base.Response{versionResult}, p.VersionPattern.FindString(versionResult.Result))

	return r
}

// GetConfig get the configuration of a source datastore from the device.
func (p *EOSCfg) GetConfig(source string) *Response {
	r := NewResponse(p.conn.Host, nil)

	// TODO make configCommandMap and then use that to fetch appropriate command based on the source
	_ = source
	configResult, err := p.conn.SendCommand("show version | i Software image version")

	if err != nil {
		return r
	}

	r.Record([]*base.Response{configResult}, p.VersionPattern.FindString(configResult.Result))

	return r
}
