package options

import (
	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util"
)

// WithStandardTransportPrivateKeyBytes sets the SSH private key and passphrase to use for SSH key based auth.
func WithStandardTransportPrivateKeyBytes(ks []byte, ps string) util.Option {
	return func(o interface{}) error {
		a, ok := o.(*transport.Standard)

		if !ok {
			return util.ErrIgnoredOption
		}

		a.PrivateKey = ks
		a.SSHArgs.PrivateKeyPassPhrase = ps

		return nil
	}
}

// WithStandardTransportExtraCiphers extends the list of ciphers supported by the standard
// transport.
func WithStandardTransportExtraCiphers(l []string) util.Option {
	return func(o interface{}) error {
		t, ok := o.(*transport.Standard)

		if !ok {
			return util.ErrIgnoredOption
		}

		t.ExtraCiphers = l

		return nil
	}
}

// WithStandardTransportExtraKexs extends the list of kext (key exchange algorithms) supported by
// the standard transport.
func WithStandardTransportExtraKexs(l []string) util.Option {
	return func(o interface{}) error {
		t, ok := o.(*transport.Standard)

		if !ok {
			return util.ErrIgnoredOption
		}

		t.ExtraKexs = l

		return nil
	}
}

// WithStandardTransportHostKeyAlgorithms sets the allowed SSH host key algorithms supported by
// the standard transport.
func WithStandardTransportHostKeyAlgorithms(algorithms []string) util.Option {
	return func(o interface{}) error {
		t, ok := o.(*transport.Standard)
		if !ok {
			return util.ErrIgnoredOption
		}

		t.HostKeyAlgorithms = algorithms

		return nil
	}
}
