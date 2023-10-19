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
	pconn *ipv4.PacketConn

	// ifps stores pointers of multicast interface.
	ifps []*net.Interface
}

type connConfig struct {
	ttl int
}

// Listen starts to receiving multicast messages.
func Listen(r *AddrResolver, opts ...ConnOption) (*Conn, error) {
	// prepare parameters.
	laddr, err := r.resolve()
	if err != nil {
		return nil, err
	}
	// configure connection
	var cfg connConfig
	for _, o := range opts {
		o.apply(&cfg)
	}
	// connect.
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return nil, err
	}
	// configure socket to use with multicast.
	pconn, iflist, err := newIPv4MulticastConn(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	// store interfaces by pointer.
	ifplist := make([]*net.Interface, 0, len(iflist))
	for i := range iflist {
		ifplist = append(ifplist, &iflist[i])
	}
	// set TTL
	if cfg.ttl > 0 {
		err := pconn.SetTTL(cfg.ttl)
		if err != nil {
			pconn.Close()
			return nil, err
		}
	}
	return &Conn{
		laddr: laddr,
		pconn: pconn,
		ifps:  ifplist,
	}, nil
}

func newIPv4MulticastConn(conn *net.UDPConn) (*ipv4.PacketConn, []net.Interface, error) {
	iflist, err := interfaces()
	if err != nil {
		return nil, nil, err
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
func joinGroupIPv4(conn *net.UDPConn, iflist []net.Interface, gaddr net.Addr) (*ipv4.PacketConn, error) {
	wrap := ipv4.NewPacketConn(conn)
	wrap.SetMulticastLoopback(true)
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
	// based net.UDPConn will be closed by mc.pconn.Close()
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
	if uaddr, ok := to.(*net.UDPAddr); ok && !uaddr.IP.IsMulticast() {
		return mc.writeToIfi(dataProv, to, nil)
	}
	sum := 0
	for _, ifp := range mc.ifps {
		n, err := mc.writeToIfi(dataProv, to, ifp)
		if err != nil {
			return 0, err
		}
		sum += n
	}
	return sum, nil
}

func (mc *Conn) writeToIfi(dataProv DataProvider, to net.Addr, ifi *net.Interface) (int, error) {
	if err := mc.pconn.SetMulticastInterface(ifi); err != nil {
		return 0, err
	}
	return mc.pconn.WriteTo(dataProv.Bytes(ifi), nil, to)
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

// ConnOption is option for Listen()
type ConnOption interface {
	apply(cfg *connConfig)
}

type connOptFunc func(*connConfig)

func (f connOptFunc) apply(cfg *connConfig) {
	f(cfg)
}

// ConnTTL returns as ConnOption that set default TTL to the connection.
func ConnTTL(ttl int) ConnOption {
	return connOptFunc(func(cfg *connConfig) {
		cfg.ttl = ttl
	})
}
