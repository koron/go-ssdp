// 存在するnet.Interfaceとそのアドレスを列挙する
package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	if err := list(); err != nil {
		log.Fatal(err)
	}
}

func list() error {
	iflist, err := net.Interfaces()
	if err != nil {
		return err
	}
	for i, nif := range iflist {
		fmt.Printf("#%d %+v\n", i, nif)
		//if !hasLinkUp(&nif) {
		//	fmt.Printf("  disabled: no link up\n")
		//	continue
		//}
		//if !hasIPv4Address(&nif) {
		//	fmt.Printf("  disabled: no IPv4 address\n")
		//	continue
		//}
		addrs, err := nif.Addrs()
		if err != nil {
			fmt.Printf("  address error: %s\n", err)
		}
		for j, a := range addrs {
			fmt.Printf("  #%d network:%s string:%s\n", j, a.Network(), a.String())
		}
		maddrs, err := nif.MulticastAddrs()
		if err != nil {
			fmt.Printf("  multicast address error: %s\n", err)
		}
		for j, a := range maddrs {
			fmt.Printf("  +%d network:%s string:%s\n", j, a.Network(), a.String())
		}
	}
	return nil
}

// hasLinkUp checks an I/F have link-up or not.
func hasLinkUp(ifi *net.Interface) bool {
	return ifi.Flags&net.FlagUp != 0
}

// hasIPv4Address checks an I/F have IPv4 address.
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
