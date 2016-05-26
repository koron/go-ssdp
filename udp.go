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

func multicastListen(localAddr string) (net.PacketConn, error) {
	conn, err := net.ListenPacket("udp4", localAddr)
	if err != nil {
		return nil, err
	}
	if err := setupMulticast(conn); err != nil {
		conn.Close()
		return nil, err
	}
	logf("listening %s for multicast", conn.LocalAddr().String())
	return conn, nil
}

func setupMulticast(conn net.PacketConn) error {
	wrap := ipv4.NewPacketConn(conn)
	wrap.SetMulticastLoopback(true)
	// add interfaces to multicast group.
	iflist, err := net.Interfaces()
	if err != nil {
		return err
	}
	empty := true
	for _, ifi := range iflist {
		if !hasRealAddress(&ifi) {
			continue
		}
		wrap.JoinGroup(&ifi, multicastAddr4)
		empty = false
	}
	if empty {
		return errors.New("no interfaces to listen")
	}
	return nil
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
