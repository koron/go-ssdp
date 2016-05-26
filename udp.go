package ssdp

import (
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

var multicastAddr4 *net.UDPAddr

func init() {
	var err error
	multicastAddr4, err = net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		panic(err)
	}
}

type packetHandler func(net.Addr, []byte) error

func readPackets(conn *net.UDPConn, timeout time.Duration, h packetHandler) error {
	buf := make([]byte, 65535)
	conn.SetReadBuffer(len(buf))
	conn.SetReadDeadline(time.Now().Add(timeout))
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				return nil
			}
			return err
		}
		if err := h(addr, buf[:n]); err != nil {
			return err
		}
	}
}

func listen(localAddr string) (*net.UDPConn, error) {
	// prepare parameters.
	laddr, err := net.ResolveUDPAddr("udp4", localAddr)
	if err != nil {
		return nil, err
	}
	// connect.
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return nil, err
	}
	// TODO: configure socket to use with multicast.
	pc := ipv4.NewPacketConn(conn)
	n, err := pc.MulticastLoopback()
	if err != nil {
		logf("MulticastLoopback() failed")
	} else {
		logf("MulticastLoopback=%t\n", n)
	}
	return conn, err
}
