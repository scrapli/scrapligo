package base

import (
	"strings"
	"time"

	"github.com/scrapli/scrapligo/channel"
)

// FullSendInteractive same as `SendInteractive` but requiring explicit options.
func (d *Driver) FullSendInteractive(
	events []*channel.SendInteractiveEvent,
	failedWhenContains []string,
	timeoutOps time.Duration,
	joinedEventInputs string,
) (*Response, error) {
	r := NewResponse(d.Host, d.Port, joinedEventInputs, failedWhenContains)

	rawResult, err := d.Channel.SendInteractive(events, timeoutOps)

	r.Record(rawResult, string(rawResult))

	if err != nil {
		r.Failed = true
	}

	return r, err
}

// SendInteractive send interactive commands to a device, accepts a slice of `SendInteractiveEvent`
// and variadic of `SendOption`s.
func (d *Driver) SendInteractive(
	events []*channel.SendInteractiveEvent,
	o ...SendOption,
) (*Response, error) {
	finalOpts := d.ParseSendOptions(o)
	joinedEventInputs := JoinEventInputs(events)

	return d.FullSendInteractive(
		events,
		finalOpts.FailedWhenContains,
		finalOpts.TimeoutOps,
		joinedEventInputs,
	)
}

// JoinEventInputs convenience function to join inputs from a `SendInteractive` method.
func JoinEventInputs(events []*channel.SendInteractiveEvent) string {
	eventInputs := make([]string, len(events))

	for _, event := range events {
		eventInputs = append(eventInputs, event.ChannelInput)
	}

	joinedEventInputs := strings.Join(eventInputs, ", ")

	return joinedEventInputs
}
