package ssdp

import (
	"errors"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

type multicastConn struct {
	laddr  *net.UDPAddr
	conn   *net.UDPConn
	iflist []net.Interface
	pconn  *ipv4.PacketConn
}

func multicastListen(localAddr string) (*multicastConn, error) {
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
	// configure socket to use with multicast.
	iflist := interfaces()
	pconn, err := joinGroupIPv4(conn, iflist, ssdpAddrIPv4)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &multicastConn{
		laddr:  laddr,
		iflist: iflist,
		conn:   conn,
		pconn:  pconn,
	}, nil
}

// joinGroupIPv4 makes the connection join to a group on interfaces.
func joinGroupIPv4(conn *net.UDPConn, iflist []net.Interface, gaddr net.Addr) (*ipv4.PacketConn, error) {
	wrap := ipv4.NewPacketConn(conn)
	wrap.SetMulticastLoopback(true)
	// add interfaces to multicast group.
	joined := 0
	for _, ifi := range iflist {
		if err := wrap.JoinGroup(&ifi, gaddr); err != nil {
			logf("failed to join group %s on %s: %s", gaddr.String(), ifi.Name, err)
			continue
		}
		joined++
		logf("joined gropup %s on %s", gaddr.String(), ifi.Name)
	}
	if joined == 0 {
		return nil, errors.New("no interfaces had joined to group")
	}
	return wrap, nil
}

func (mc *multicastConn) Close() error {
	if err := mc.pconn.Close(); err != nil {
		return err
	}
	return nil
}

func (mc *multicastConn) WriteTo(data []byte, to net.Addr) (int, error) {
	for _, ifi := range mc.iflist {
		if err := mc.pconn.SetMulticastInterface(&ifi); err != nil {
			logf("HERE1: %s", err)
			return 0, err
		}
		if _, err := mc.pconn.WriteTo(data, nil, to); err != nil {
			logf("HERE2: %s", err)
			return 0, err
		}
	}
	return len(data), nil
}

func (mc *multicastConn) LocalAddr() net.Addr {
	return mc.laddr
}

func (mc *multicastConn) readPackets(timeout time.Duration, h packetHandler) error {
	buf := make([]byte, 65535)
	if timeout > 0 {
		mc.pconn.SetReadDeadline(time.Now().Add(timeout))
	}
	for {
		n, _, addr, err := mc.pconn.ReadFrom(buf)
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
