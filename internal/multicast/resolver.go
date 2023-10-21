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

func (ar *addrResolver) Resolve() (*net.UDPAddr, error) {
	ar.once.Do(func() {
		ar.udp, ar.err = net.ResolveUDPAddr("udp4", ar.addr)
	})
	return ar.udp, ar.err
}

// NewResolver creates a resolver from string.
func NewResolver(addr string) Resolver {
	return &addrResolver{addr: addr}
}
