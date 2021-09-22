package netconf

import (
	"errors"
	"regexp"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/transport"
)

// ErrCapabilitiesExchangeFailed error for failure of capabilities exchange.
var ErrCapabilitiesExchangeFailed = errors.New("failure during capabilities exchange")

// Driver the Netconf Driver struct, extends the base driver, and embeds netconf channel.
type Driver struct {
	base.Driver
	NetconfChannel      *Channel
	readableDatastores  [][]byte
	writeableDatastores [][]byte
	messageID           int
	// these are not settable via NewNetconfDriver, just set them manually before open if you care
	StripNamespaces  bool
	StrictDatastores bool
}

// NewNetconfDriver return an instance of netconf `Driver`.
func NewNetconfDriver(
	host string,
	options ...base.Option,
) (*Driver, error) {
	options = append([]base.Option{base.WithPort(DefaultPort)}, options...)

	newDriver, err := base.NewDriver(host, options...)

	if err != nil {
		return nil, err
	}

	d := &Driver{
		Driver:           *newDriver,
		StripNamespaces:  false,
		StrictDatastores: true,
		messageID:        101,
	}

	nc := &Channel{
		BaseChannel: d.Channel,
		serverEcho:  d.NetconfEcho,
	}

	d.NetconfChannel = nc

	// ignoring user input on comms prompt pattern too as we know what we need to look for
	d.Channel.CommsPromptPattern = regexp.MustCompile(Version10DelimiterPattern)

	// temp for appeasing linting until they are used
	_ = d.readableDatastores
	_ = d.writeableDatastores

	return d, nil
}

// Open netconf open method, opens transport in netconf "mode".
func (d *Driver) Open() error {
	err := d.Transport.OpenNetconf()
	if err != nil {
		return err
	}

	var authenticationBuf []byte

	if d.TransportType == transport.SystemTransportName {
		r, authErr := d.Channel.AuthenticateSSH(d.AuthPassword, d.AuthPrivateKeyPassphrase)
		if authErr != nil {
			logging.LogError(
				d.FormatLogMessage("error", "authentication failed, connection not opened"),
			)

			return authErr
		}

		authenticationBuf = r
	}

	err = d.NetconfChannel.OpenNetconf(authenticationBuf)

	return err
}

// Close closes the connection.
func (d *Driver) Close() error {
	err := d.Transport.Close()
	return err
}

// ParseNetconfOptions parse provided netconf options.
func (d *Driver) ParseNetconfOptions(
	o []Option,
) *Options {
	finalOpts := &Options{
		Filter:      DefaultNetconfOptionsFilter,
		FilterType:  DefaultNetconfOptionsFilterType,
		DefaultType: DefaultNetconfOptionsDefaultType,
	}

	if len(o) > 0 && o[0] != nil {
		for _, option := range o {
			option(finalOpts)
		}
	}

	return finalOpts
}

func (d *Driver) finalizeAndSendMessage(
	netconfMessage *Message,
) (*Response, error) {
	bytesNetconfMessage, err := d.NetconfChannel.BuildFinalMessage(netconfMessage)
	if err != nil {
		return NewNetconfResponse(
			d.Host,
			d.NetconfChannel.SelectedNetconfVersion,
			d.Transport.BaseTransportArgs.Port,
			[]byte(""),
			netconfMessage,
			d.StripNamespaces,
		), err
	}

	r := NewNetconfResponse(
		d.Host,
		d.NetconfChannel.SelectedNetconfVersion,
		d.Transport.BaseTransportArgs.Port,
		bytesNetconfMessage,
		netconfMessage,
		d.StripNamespaces,
	)

	channelResponse, err := d.NetconfChannel.SendInputNetconf(bytesNetconfMessage)

	r.Record(channelResponse)

	if err != nil {
		return r, err
	}

	return r, nil
}
