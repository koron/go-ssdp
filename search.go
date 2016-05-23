package ssdp

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

// Service is discovered service.
type Service struct {
	// TODO:
}

// Search searchs services by SSDP.
func Search(searchType string, waitTime int) ([]Service, error) {
	// prepare parameters.
	laddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:0")
	if err != nil {
		return nil, err
	}

	// connect.
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	logf("search on %s", conn.LocalAddr().String())

	// send request.
	msg, err := buildSearch(multicastAddr4, searchType, waitTime)
	if err != nil {
		return nil, err
	}
	if _, err := conn.WriteTo(msg, multicastAddr4); err != nil {
		return nil, err
	}

	// wait response.
	var list []Service
	h := func(a net.Addr, d []byte) error {
		srv, err := parseService(a, d)
		if err != nil {
			logf("invalid search response: %s", err)
			return nil
		}
		list = append(list, *srv)
		return nil
	}
	if err := readPackets(conn, time.Duration(waitTime), h); err != nil {
		return nil, err
	}
	return list, err
}

func buildSearch(raddr net.Addr, searchType string, waitTime int) ([]byte, error) {
	b := new(bytes.Buffer)
	// FIXME: error should be checked.
	b.WriteString("M-SEARCH * HTTP/1.1\r\n")
	fmt.Fprintf(b, "HOST: %s\r\n", raddr.String())
	fmt.Fprintf(b, "MAN: %q\r\n", "ssdp:discover")
	fmt.Fprintf(b, "MX: %d\r\n", waitTime)
	fmt.Fprintf(b, "ST: %s\r\n", searchType)
	b.WriteString("\r\n")
	return b.Bytes(), nil
}

func parseService(addr net.Addr, data []byte) (*Service, error) {
	// TODO:
	return nil, nil
}
