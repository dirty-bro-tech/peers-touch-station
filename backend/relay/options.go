package relay

import (
	"context"
	"fmt"

	"github.com/dirty-bro-tech/peers-touch-go/core/server"
)

type optionsKey struct{}

type HeadProtocol string

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

type Option func(*Options)

func KeyFile(keyFile string) server.SubServerOption {
	return func(o *server.SubServerOptions) {
		optionWrap(o, func(opts *Options) {
			opts.KeyFile = keyFile
		})
	}
}

func Addresses(adds ...Addr) server.SubServerOption {
	return func(o *server.SubServerOptions) {
		optionWrap(o, func(opts *Options) {
			opts.Addresses = append(opts.Addresses, adds...)
		})
	}
}

type GetOptions struct {
}

type GetOption func(*GetOptions)

func optionWrap(o *server.SubServerOptions, f func(*Options)) {
	if o.Ctx == nil {
		o.Ctx = context.Background()
	}

	var opts *Options
	if o.Ctx.Value(optionsKey{}) == nil {
		opts = &Options{}
		o.Ctx = context.WithValue(o.Ctx, optionsKey{}, opts)
	} else {
		opts = o.Ctx.Value(optionsKey{}).(*Options)
	}

	f(opts)
}
