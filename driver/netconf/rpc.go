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

func forceSelfClosingTags(b []byte) []byte {
	ncPatterns := getNetconfPatterns()

	emptyTagIdxs := ncPatterns.emptyTags.FindAllSubmatchIndex(b, -1)

	var nb []byte

	for _, idx := range emptyTagIdxs {
		// get everything in b up till the first of the submatch indexes (this is the start of an
		// "empty" <thing></thing> tag), then get the name of the tag and put it in a self-closing
		// tag.
		nb = append(b[0:idx[0]], fmt.Sprintf("<%s/>", b[idx[2]:idx[3]])...) //nolint: gocritic

		// finally, append everything *after* the submatch indexes
		nb = append(nb, b[len(b)-(len(b)-idx[1]):]...)
	}

	return nb
}

func (d *Driver) sendRPC(
	m *message,
	op *OperationOptions,
) (*response.NetconfResponse, error) {
	b, err := m.serialize(d.SelectedVersion)
	if err != nil {
		return nil, err
	}

	if d.ForceSelfClosingTags {
		b = forceSelfClosingTags(b)
	}

	r := response.NewNetconfResponse(
		b,
		d.Transport.GetHost(),
		d.Transport.GetPort(),
		d.SelectedVersion,
	)

	d.storeMessage(m.MessageID, nil)

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

	done := make(chan bool)

	go func() {
		for d.getMessage(m.MessageID) == nil {
			time.Sleep(d.Channel.ReadDelay)
		}

		done <- true
	}()

	timer := time.NewTimer(d.Channel.GetTimeout(op.Timeout))

	select {
	case err = <-d.errs:
		return nil, err
	case <-timer.C:
		d.Logger.Critical("channel timeout sending input to device")

		return nil, fmt.Errorf("%w: channel timeout sending input to device", util.ErrTimeoutError)
	case <-done:
	}

	r.Record(d.getMessage(m.MessageID))

	return r, nil
}
