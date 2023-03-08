package netconf

import (
	"fmt"
	"time"

	"github.com/scrapli/scrapligo/response"
	"github.com/scrapli/scrapligo/util"
)

func (d *Driver) buildRPCElem(
	filter string,
) *message {
	netconfInput := d.buildPayload(filter)

	return netconfInput
}

// RPC executes a "bare" RPC against the NETCONF server.
func (d *Driver) RPC(opts ...util.Option) (*response.NetconfResponse, error) {
	op, err := NewOperation(opts...)
	if err != nil {
		return nil, err
	}

	return d.sendRPC(d.buildRPCElem(op.Filter), op)
}

// ForceSelfClosingTags accepts a netconf looking xml byte slice and forces any "empty" tags (tags
// without attributes) to use self-closing tags. For example:
// `<running> </running>`
// Would be converted to:
// `<running/>`.
func ForceSelfClosingTags(b []byte) []byte {
	ncPatterns := getNetconfPatterns()

	r := ncPatterns.emptyTags.ReplaceAll(b, []byte("<$1$2/>"))

	return r
}

func (d *Driver) sendRPC(
	m *message,
	op *OperationOptions,
) (*response.NetconfResponse, error) {
	if d.ForceSelfClosingTags {
		d.Logger.Debug("ForceSelfClosingTags is true, enforcing...")
	}

	b, err := m.serialize(d.SelectedVersion, d.ForceSelfClosingTags)
	if err != nil {
		return nil, err
	}

	d.Logger.Debugf("sending finalized rpc payload:\n%s", string(b))

	r := response.NewNetconfResponse(
		b,
		d.Transport.GetHost(),
		d.Transport.GetPort(),
		d.SelectedVersion,
	)

	err = d.Channel.WriteAndReturn(b, false)
	if err != nil {
		return nil, err
	}

	if d.SelectedVersion == V1Dot1 {
		err = d.Channel.WriteReturn()
		if err != nil {
			return nil, err
		}
	}

	done := make(chan []byte)

	go func() {
		var data []byte

		for {
			data = d.getMessage(m.MessageID)
			if data != nil {
				break
			}

			time.Sleep(5 * time.Microsecond) //nolint: gomnd
		}

		done <- data
	}()

	timer := time.NewTimer(d.Channel.GetTimeout(op.Timeout))

	select {
	case err = <-d.errs:
		return nil, err
	case <-timer.C:
		d.Logger.Critical("channel timeout sending input to device")

		return nil, fmt.Errorf("%w: channel timeout sending input to device", util.ErrTimeoutError)
	case data := <-done:
		r.Record(data)
	}

	return r, nil
}
