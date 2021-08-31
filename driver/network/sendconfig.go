package network

import (
	"strings"
	"time"

	"github.com/scrapli/scrapligo/driver/base"
)

// SendConfig send configuration string to the device.
func (d *Driver) SendConfig(c string, o ...base.SendOption) (*base.Response, error) {
	sc := strings.Split(c, "\n")
	m, err := d.SendConfigs(sc, o...)

	r := base.NewResponse(
		d.Host,
		d.Transport.BaseTransportArgs.Port,
		c,
		m.Responses[0].FailedWhenContains,
	)

	if err != nil {
		return r, err
	}

	individualResponses := make([]string, len(sc))
	for _, response := range m.Responses {
		individualResponses = append(individualResponses, response.Result)
	}

	r.StartTime = m.StartTime
	r.EndTime = time.Now()
	r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()
	r.Result = strings.Join(individualResponses, "\n")

	return r, nil
}
