package multicast

import (
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"github.com/koron/go-ssdp/internal/ssdplog"
	"golang.org/x/net/ipv4"
)

// Conn is multicast connection.
type Conn struct {
	laddr *net.UDPAddr
	conn  *net.UDPConn
	pconn *ipv4.PacketConn

	// to use system assigned multicast interface, set this empty or nil.
	iflist []net.Interface
}

// Listen starts to receiving multicast messages.
func Listen(r *AddrResolver) (*Conn, error) {
	return listen(r, SystemAssignedInterface)
}

// ListenInterfaces starts to receiving multicast messages on all provided interfaces.
// See InterfaceProvider also.
func ListenInterfaces(r *AddrResolver) (*Conn, error) {
	return listen(r, false)
}

// ListenSystem starts to receiving multicast messages on system assigned multicast interface.
// InterfacesProvider will be ignored.
func ListenSystem(r *AddrResolver) (*Conn, error) {
	return listen(r, true)
}

// listen starts to receiving multicast messages.
// Use the system assigned multicast interface when sysIf is true.
func listen(r *AddrResolver, sysIf bool) (*Conn, error) {
	// prepare parameters.
	laddr, err := r.resolve()
	if err != nil {
		return nil, err
	}
	// connect.
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return nil, err
	}
	// configure socket to use with multicast.
	pconn, iflist, err := newIPv4MulticastConn(conn, sysIf)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &Conn{
		laddr:  laddr,
		conn:   conn,
		pconn:  pconn,
		iflist: iflist,
	}, nil
}

// newIPv4MulticastConn create a new multicast connection.
// 2nd return parameter will be nil when sysIf is true.
func newIPv4MulticastConn(conn *net.UDPConn, sysIf bool) (*ipv4.PacketConn, []net.Interface, error) {
	// sysIf: use system assigned multicast interface.
	// the empty iflist indicate it.
	var iflist []net.Interface
	if !sysIf {
		list, err := interfaces()
		if err != nil {
			return nil, nil, err
		}
		iflist = list
	}
	addr, err := SendAddr()
	if err != nil {
		return nil, nil, err
	}
	pconn, err := joinGroupIPv4(conn, iflist, addr)
	if err != nil {
		return nil, nil, err
	}
	return pconn, iflist, nil
}

// joinGroupIPv4 makes the connection join to a group on interfaces.
// This trys to use system assigned when iflist is nil or empty.
func joinGroupIPv4(conn *net.UDPConn, iflist []net.Interface, gaddr net.Addr) (*ipv4.PacketConn, error) {
	wrap := ipv4.NewPacketConn(conn)
	wrap.SetMulticastLoopback(true)

	// try to use the system assigned multicast interface when iflist is empty.
	if len(iflist) == 0 {
		if err := wrap.JoinGroup(nil, gaddr); err != nil {
			ssdplog.Printf("failed to join group %s on system assigned multicast interface: %s", gaddr.String(), err)
			return nil, errors.New("no system assigned multicast interfaces had joined to group")
		}
		ssdplog.Printf("joined group %s on system assigned multicast interface", gaddr.String())
		return wrap, nil
	}

	// add interfaces to multicast group.
	joined := 0
	for _, ifi := range iflist {
		if err := wrap.JoinGroup(&ifi, gaddr); err != nil {
			ssdplog.Printf("failed to join group %s on %s: %s", gaddr.String(), ifi.Name, err)
			continue
		}
		joined++
		ssdplog.Printf("joined group %s on %s (#%d)", gaddr.String(), ifi.Name, ifi.Index)
	}
	if joined == 0 {
		return nil, errors.New("no interfaces had joined to group")
	}
	return wrap, nil
}

// Close closes a multicast connection.
func (mc *Conn) Close() error {
	if err := mc.pconn.Close(); err != nil {
		return err
	}
	// mc.conn is closed by mc.pconn.Close()
	return nil
}

// DataProvider provides a body of multicast message to send.
type DataProvider interface {
	Bytes(*net.Interface) []byte
}

//type multicastDataProviderFunc func(*net.Interface) []byte
//
//func (f multicastDataProviderFunc) Bytes(ifi *net.Interface) []byte {
//	return f(ifi)
//}

type BytesDataProvider []byte

func (b BytesDataProvider) Bytes(ifi *net.Interface) []byte {
	return []byte(b)
}

// WriteTo sends a multicast message to interfaces.
func (mc *Conn) WriteTo(dataProv DataProvider, to net.Addr) (int, error) {
	// Send a multicast message directory when recipient "to" address is not multicast.
	if uaddr, ok := to.(*net.UDPAddr); ok && !uaddr.IP.IsMulticast() {
		return mc.conn.WriteTo(dataProv.Bytes(nil), to)
	}
	// Send a multicast message directory when using system assigned multicast interface.
	if len(mc.iflist) == 0 {
		return mc.conn.WriteTo(dataProv.Bytes(nil), to)
	}
	// Send a multicast message to all interfaces (iflist).
	sum := 0
	for _, ifi := range mc.iflist {
		if err := mc.pconn.SetMulticastInterface(&ifi); err != nil {
			return 0, err
		}
		n, err := mc.pconn.WriteTo(dataProv.Bytes(&ifi), nil, to)
		if err != nil {
			return 0, err
		}
		sum += n
	}
	return sum, nil
}

// LocalAddr returns local address to listen multicast packets.
func (mc *Conn) LocalAddr() net.Addr {
	return mc.laddr
}

// ReadPackets reads multicast packets.
func (mc *Conn) ReadPackets(timeout time.Duration, h PacketHandler) error {
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
			if strings.Contains(err.Error(), "use of closed network connection") {
				return io.EOF
			}
			return err
		}
		if err := h(addr, buf[:n]); err != nil {
			return err
		}
	}
}
