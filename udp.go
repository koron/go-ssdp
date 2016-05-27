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
	// FIXME: configure socket to use with multicast.
	return conn, err
}

func multicastListen(localAddr string, iflist []net.Interface) (net.PacketConn, error) {
	conn, err := net.ListenPacket("udp4", localAddr)
	if err != nil {
		return nil, err
	}
	if err := joinGroupIPv4(conn, iflist, multicastAddr4); err != nil {
		conn.Close()
		return nil, err
	}
	logf("listening %s for multicast", conn.LocalAddr().String())
	return conn, nil
}

// joinGroupIPv4 makes the connection join to a group on interfaces.
func joinGroupIPv4(conn net.PacketConn, iflist []net.Interface, gaddr net.Addr) error {
	wrap := ipv4.NewPacketConn(conn)
	wrap.SetMulticastLoopback(true)
	if len(iflist) == 0 {
		iflist = interfacesIPv4()
	}
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
		return errors.New("no interfaces had joined to group")
	}
	return nil
}

func interfacesIPv4() []net.Interface {
	iflist, err := net.Interfaces()
	if err != nil {
		return nil
	}
	list := make([]net.Interface, 0, len(iflist))
	for _, ifi := range iflist {
		if !hasIPv4Address(&ifi) {
			continue
		}
		list = append(list, ifi)
	}
	return list
}

func hasIPv4Address(ifi *net.Interface) bool {
	addrs, err := ifi.Addrs()
	if err != nil {
		return false
	}
	for _, a := range addrs {
		ip, _, err := net.ParseCIDR(a.String())
		if err != nil {
			continue
		}
		if len(ip.To4()) == net.IPv4len && !ip.IsUnspecified() {
			return true
		}
	}
	return false
}
