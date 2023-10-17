package ssdp

import (
	"log"
	"net"

	"github.com/koron/go-ssdp/internal/multicast"
	"github.com/koron/go-ssdp/internal/ssdplog"
)

func init() {
	multicast.InterfacesProvider = func() []net.Interface {
		return Interfaces
	}
	ssdplog.LoggerProvider = func() *log.Logger {
		return Logger
	}
}

// Interfaces specify target interfaces to multicast.  If no interfaces are
// specified, all interfaces will be used.
var Interfaces []net.Interface

// Logger is default logger for SSDP module.
var Logger *log.Logger

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

func defaultConnOpts() []multicast.ConnOption {
	var opts []multicast.ConnOption
	if multicastTTL > 0 {
		opts = append(opts, multicast.ConnTTL(multicastTTL))
	}
	if multicastSystemAssignedInterface {
		opts = append(opts, multicast.ConnSystemAssginedInterface())
	}
	return opts
}

var multicastTTL int

// SetMulticastTTL sets default TTL of SSDP's UDP packets.
// 0 default, 1 or greater set TTL.
func SetMulticastTTL(ttl int) {
	multicastTTL = ttl
}

var multicastSystemAssignedInterface bool

// SetMulticastSystemAssignedInterface updates state whether using the system
// assigned multicast interface or provided interfaces.
// Default is "false", it uses provided interface (see Interfaces also).
func SetMulticastSystemAssignedInterface(enable bool) {
	multicastSystemAssignedInterface = enable
}

// GetMulticastSystemAssignedInterface returns state using the system assigned
// multicast interface or provided interfaces.
func GetMulticastSystemAssignedInterface() bool {
	return multicastSystemAssignedInterface
}
