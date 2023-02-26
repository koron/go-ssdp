// IPv6 マルチキャストを受信する
//
// Windows では動かなかったがFreeBSDでは動いた
package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

func main() {
	if err := listen(); err != nil {
		log.Fatal(err)
	}
}

//const addrStr = "239.255.255.250:1900"
const addrStr = "[FF02::C]:1900"

func listen() error {
	ifi, err := net.InterfaceByIndex(40)
	if err != nil {
		return err
	}
	fmt.Printf("interface: %+v\n", ifi)
	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		return err
	}
	conn, err := net.ListenMulticastUDP("udp", ifi, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	buf := make([]byte, 65535)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				break
			}
			return err
		}
		raw := buf[:n]
		if bytes.HasPrefix(raw, []byte("M-SEARCH ")) {
			fmt.Printf("detect M-SEARCH from %s\n", addr)
			continue
		}
		if bytes.HasPrefix(raw, []byte("NOTIFY ")) {
			fmt.Printf("detect NOTIFY from %s\n", addr)
			continue
		}
		x := bytes.Index(raw, []byte("\r\n"))
		fmt.Printf("unexpected method from %s: %q", addr, string(raw[:x]))
	}
	return nil
}
