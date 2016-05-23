package ssdp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Service is discovered service.
type Service struct {
	SearchType string
	USN        string
	Location   string
	Server     string
	RawHeader  http.Header

	maxAge *int
}

func (s *Service) MaxAge() int {
	if s.maxAge != nil {
		return *s.maxAge
	}
	s.maxAge = new(int)
	*s.maxAge = -1
	if v := s.RawHeader.Get("CACHE-CONTROL"); v != "" {
	}
	// TODO: parse CACHE-CONTROL
	return *s.maxAge
}

// Search searchs services by SSDP.
func Search(searchType string, waitSec int) ([]Service, error) {
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
	msg, err := buildSearch(multicastAddr4, searchType, waitSec)
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
	if err := readPackets(conn, time.Duration(waitSec), h); err != nil {
		return nil, err
	}

	return list, err
}

func buildSearch(raddr net.Addr, searchType string, waitSec int) ([]byte, error) {
	b := new(bytes.Buffer)
	// FIXME: error should be checked.
	b.WriteString("M-SEARCH * HTTP/1.1\r\n")
	fmt.Fprintf(b, "HOST: %s\r\n", raddr.String())
	fmt.Fprintf(b, "MAN: %q\r\n", "ssdp:discover")
	fmt.Fprintf(b, "MX: %d\r\n", waitSec)
	fmt.Fprintf(b, "ST: %s\r\n", searchType)
	b.WriteString("\r\n")
	return b.Bytes(), nil
}

var (
	errWithoutHTTPPrefix = errors.New("without HTTP prefix")
)

func parseService(addr net.Addr, data []byte) (*Service, error) {
	if !bytes.HasPrefix(data, []byte("HTTP")) {
		return nil, errWithoutHTTPPrefix
	}
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(data)), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &Service{
		SearchType: resp.Header.Get("ST"),
		USN:        resp.Header.Get("USN"),
		Location:   resp.Header.Get("LOCATION"),
		Server:     resp.Header.Get("SERVER"),
		RawHeader:  resp.Header,
	}, nil
}
