package ssdp

import (
	"net"

	"github.com/koron/go-ssdp/internal/multicast"
)

func init() {
	multicast.InterfacesProvider = func() []net.Interface {
		return Interfaces
	}
}

// Interfaces specify target interfaces to multicast.  If no interfaces are
// specified, all interfaces will be used.
var Interfaces []net.Interface

// SetMulticastRecvAddrIPv4 updates multicast address where to receive packets.
// This never fail now.
func SetMulticastRecvAddrIPv4(addr string) error {
	return multicast.SetRecvAddrIPv4(addr)
}

// SetMulticastSendAddrIPv4 updates a UDP address to send multicast packets.
// This never fail now.
func SetMulticastSendAddrIPv4(addr string) error {
	return multicast.SetSendAddrIPv4(addr)
}
