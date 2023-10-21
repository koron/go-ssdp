package multicast

import (
	"net"
	"sync"
)

type PacketHandler func(net.Addr, []byte) error

type AddrResolver struct {
	Addr string

	mu  sync.RWMutex
	udp *net.UDPAddr
	err error
}

func (r *AddrResolver) setAddress(addr string) {
	r.mu.Lock()
	r.Addr = addr
	r.udp = nil
	r.err = nil
	r.mu.Unlock()
}

func (r *AddrResolver) resolve() (*net.UDPAddr, error) {
	r.mu.RLock()
	if err := r.err; err != nil {
		r.mu.RUnlock()
		return nil, err
	}
	if udp := r.udp; udp != nil {
		r.mu.RUnlock()
		return udp, nil
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()
	r.udp, r.err = net.ResolveUDPAddr("udp4", r.Addr)
	return r.udp, r.err
}

var sendAddrResolver = &AddrResolver{Addr: "239.255.255.250:1900"}

// SendAddr returns a remote address for multicast UDP.
func SendAddr0() (*net.UDPAddr, error) {
	return sendAddrResolver.resolve()
}
