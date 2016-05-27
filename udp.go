package ssdp

import (
	"errors"
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

func multicastDial(localAddr string) (*net.UDPConn, error) {
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

func multicastListen(localAddr string, iflist []net.Interface) (net.PacketConn, error) {
	conn, err := net.ListenPacket("udp4", localAddr)
	if err != nil {
		return nil, err
	}
	if err := joinGroup(conn, iflist, multicastAddr4); err != nil {
		conn.Close()
		return nil, err
	}
	logf("listening %s for multicast", conn.LocalAddr().String())
	return conn, nil
}

// joinGroup makes the connection join to a group on interfaces.
// when iflist is empty, it will use default interfaces provided by
// defaultInterfaces().
func joinGroup(conn net.PacketConn, iflist []net.Interface, gaddr net.Addr) error {
	wrap := ipv4.NewPacketConn(conn)
	wrap.SetMulticastLoopback(true)
	if len(iflist) == 0 {
		iflist = defaultInterfaces()
		if len(iflist) == 0 {
			return errors.New("no interfaces to join group")
		}
	}
	// add interfaces to multicast group.
	for _, ifi := range iflist {
		if err := wrap.JoinGroup(&ifi, gaddr); err != nil {
			logf("failed to join group %s on %s: %s", gaddr.String(), ifi.Name, err)
		}
	}
	return nil
}

func defaultInterfaces() []net.Interface {
	iflist, err := net.Interfaces()
	if err != nil {
		return nil
	}
	list := make([]net.Interface, 0, len(iflist))
	for _, ifi := range iflist {
		if !hasRealAddress(&ifi) {
			continue
		}
		list = append(list, ifi)
	}
	return list
}

func hasRealAddress(ifi *net.Interface) bool {
	addrs, err := ifi.Addrs()
	if err != nil {
		return false
	}
	for _, a := range addrs {
		ip := net.ParseIP(a.String())
		if !ip.IsUnspecified() {
			return true
		}
	}
	return false
}
