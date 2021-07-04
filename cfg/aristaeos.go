package cfg

import (
	"regexp"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

type EOSCfg struct {
	conn           *network.Driver
	versionPattern *regexp.Regexp
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
		versionPattern: regexp.MustCompile(`(?i)\d+\.\d+\.[a-z0-9\-]+(\.\d+[a-z]?)?`),
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

	r.Record([]*base.Response{versionResult}, p.versionPattern.FindString(versionResult.Result))

	return r
}
