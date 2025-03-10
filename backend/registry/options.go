package registry

import "fmt"

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
	Addresses addresses
	KeyFile   string
}

type Option func(*Options)

func KeyFile(keyFile string) Option {
	return func(o *Options) {
		o.KeyFile = keyFile
	}
}

func Addresses(adds ...Addr) Option {
	return func(o *Options) {
		o.Addresses = append(o.Addresses, adds...)
	}
}

type GetOptions struct {
}

type GetOption func(*GetOptions)
