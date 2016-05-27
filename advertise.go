package ssdp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
)

type message struct {
	to   net.Addr
	data []byte
}

// Advertiser is a server to advertise a service.
type Advertiser struct {
	st       string
	usn      string
	location string
	server   string
	maxAge   int

	conn net.PacketConn
	ch   chan *message
	wg   sync.WaitGroup
}

// Advertise starts advertisement of service.
func Advertise(st, usn, location, server string, maxAge int, iflist []net.Interface) (*Advertiser, error) {
	conn, err := multicastListen("0.0.0.0:1900", iflist)
	if err != nil {
		return nil, err
	}
	a := &Advertiser{
		st:       st,
		usn:      usn,
		location: location,
		server:   server,
		maxAge:   maxAge,
		conn:     conn,
		ch:       make(chan *message),
	}
	a.wg.Add(2)
	go func() {
		a.sendMain()
		a.wg.Done()
	}()
	go func() {
		a.serve()
		a.wg.Done()
	}()
	return a, nil
}

func (a *Advertiser) serve() error {
	buf := make([]byte, 65535)
	for {
		n, addr, err := a.conn.ReadFrom(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		msg := buf[:n]
		if err := a.handleRaw(addr, msg); err != nil {
			logf("failed to handle message: %s", err)
		}
	}
}

func (a *Advertiser) sendMain() error {
	for {
		select {
		case msg, ok := <-a.ch:
			if !ok {
				return nil
			}
			_, err := a.conn.WriteTo(msg.data, msg.to)
			if err != nil {
				if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
					logf("failed to send: %s", err)
				}
			}
		}
	}
}

func (a *Advertiser) handleRaw(from net.Addr, raw []byte) error {
	if !bytes.HasPrefix(raw, []byte("M-SEARCH ")) {
		// unexpected method.
		return nil
	}
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(raw)))
	if err != nil {
		return err
	}
	if man := req.Header.Get("MAN"); man != `"ssdp:discover"` {
		return fmt.Errorf("unexpected MAN: %s", man)
	}
	st := req.Header.Get("ST")
	if st != All && st != RootDevice && st != a.st {
		// skip when ST is not matched/expected.
		return nil
	}
	// build and send a response.
	msg, err := buildOK(a.st, a.usn, a.location, a.server, a.maxAge)
	if err != nil {
		return err
	}
	a.ch <- &message{to: from, data: msg}
	return nil
}

func buildOK(st, usn, location, server string, maxAge int) ([]byte, error) {
	b := new(bytes.Buffer)
	// FIXME: error should be checked.
	b.WriteString("HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(b, "ST: %s\r\n", st)
	fmt.Fprintf(b, "USN: %s\r\n", usn)
	if location != "" {
		fmt.Fprintf(b, "LOCATION: %s\r\n", location)
	}
	if server != "" {
		fmt.Fprintf(b, "SERVER: %s\r\n", server)
	}
	fmt.Fprintf(b, "CACHE-CONTROL: max-age=%d\r\n", maxAge)
	b.WriteString("\r\n")
	return b.Bytes(), nil
}

// Close stops advertisement.
func (a *Advertiser) Close() error {
	if a.conn != nil {
		close(a.ch)
		a.conn.Close()
		a.conn = nil
		a.wg.Wait()
	}
	return nil
}

// Alive announces ssdp:alive message.
func (a *Advertiser) Alive() error {
	msg, err := buildAlive(multicastAddr4, a.st, a.usn, a.location, a.server,
		a.maxAge)
	if err != nil {
		return err
	}
	a.ch <- &message{to: multicastAddr4, data: msg}
	return nil
}

// Bye announces ssdp:byebye message.
func (a *Advertiser) Bye() error {
	msg, err := buildBye(multicastAddr4, a.st, a.usn)
	if err != nil {
		return err
	}
	a.ch <- &message{to: multicastAddr4, data: msg}
	return nil
}
