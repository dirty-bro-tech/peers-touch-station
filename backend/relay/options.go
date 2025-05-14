package relay

import (
	"fmt"
	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
)

type optionsKey struct{}

type HeadProtocol string

var wrapper = option.NewWrapper[Options](optionsKey{}, func(options *option.Options) *Options {
	return BootstrapOptions()
})

const (
	HeadProtocolIP4 HeadProtocol = "ip4"
	HeadProtocolIP6 HeadProtocol = "ip6"
)

type TransportProtocol string

const (
	TransportProtocolTCP  TransportProtocol = "tcp"
	TransportProtocolUDP  TransportProtocol = "udp"
	TransportProtocolQUIC TransportProtocol = "quic"
)

type Addr struct {
	HeadProtocol      HeadProtocol
	Address           string
	TransportProtocol TransportProtocol
	Port              int
}

func (a Addr) String() string {
	return fmt.Sprintf("/%s/%s/%s/%d", a.HeadProtocol, a.Address, a.TransportProtocol, a.Port)
}

type addresses []Addr

func (a addresses) String() []string {
	ret := make([]string, len(a))
	for i, addr := range a {
		ret[i] = addr.String()
	}

	return ret
}

type Options struct {
	*server.SubServerOptions

	Addresses addresses
	KeyFile   string
}

func KeyFile(keyFile string) *option.Option {
	return wrapper.Wrap(func(opts *Options) {
		opts.KeyFile = keyFile
	})
}

func Addresses(adds ...Addr) *option.Option {
	return wrapper.Wrap(func(opts *Options) {
		opts.Addresses = append(opts.Addresses, adds...)
	})
}

type GetOptions struct {
}

type GetOption func(*GetOptions)
