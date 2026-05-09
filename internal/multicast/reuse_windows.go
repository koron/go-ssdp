//go:build windows

package multicast

import (
	"context"
	"net"
	"syscall"

	"golang.org/x/sys/windows"
)

func control(network, address string, c syscall.RawConn) error {
	var err error
	c.Control(func(fd uintptr) {
		err = windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1)
	})
	return err
}

func listenUDP(ctx context.Context, laddr *net.UDPAddr) (*net.UDPConn, error) {
	lc := net.ListenConfig{
		Control: control,
	}
	lp, err := lc.ListenPacket(ctx, "udp4", laddr.String())
	if err != nil {
		return nil, err
	}
	return lp.(*net.UDPConn), nil
}
