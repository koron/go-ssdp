//go:build !windows

package multicast

import (
	"context"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func control(network, address string, c syscall.RawConn) error {
	var err error
	c.Control(func(fd uintptr) {
		err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
		if err != nil {
			return
		}
		err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
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
