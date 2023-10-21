package multicast

import (
	"net"
	"sync"
)

// Resolver resolves *net.UDPAddr from something.
type Resolver interface {
	Resolve() (*net.UDPAddr, error)
}

type addrResolver struct {
	once sync.Once
	addr string
	udp  *net.UDPAddr
	err  error
}

func newAddrResolver(addr string) *addrResolver {
	return &addrResolver{addr: addr}
}

func (ar *addrResolver) Resolve() (*net.UDPAddr, error) {
	ar.once.Do(func() {
		ar.udp, ar.err = net.ResolveUDPAddr("udp4", ar.addr)
	})
	return ar.udp, ar.err
}

// LocalAddrResolver is a local address resolver for multicast UDP
var LocalAddrResolver Resolver = newAddrResolver("224.0.0.1:1900")

// RemoteAddrResolver is a remote address resolver for multicast UDP
var RemoteAddrResolver Resolver = newAddrResolver("239.255.255.250:1900")

// AddressResolver creates a resolver from string.
func AddressResolver(addr string) Resolver {
	return newAddrResolver(addr)
}
