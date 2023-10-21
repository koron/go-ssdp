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
