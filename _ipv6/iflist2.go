// uplinkがありmulticast可能なnet.Interfaceを列挙する
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
	fmt.Println("original:")
	for i, nif := range iflist {
		fmt.Printf("  #%d %+v\n", i, nif)
	}
	fmt.Println("filtered:")
	for i, nif := range iflist {
		if !hasLinkUp(&nif) || !hasMulticast(&nif) || !hasIPv4Address(&nif) {
			continue
		}
		fmt.Printf("  #%d %+v\n", i, nif)
	}
	return nil
}

// hasLinkUp checks an I/F have link-up or not.
func hasLinkUp(ifi *net.Interface) bool {
	return ifi.Flags&net.FlagUp != 0
}

// hasMulticast checks an I/F supports multicast or not.
func hasMulticast(ifi *net.Interface) bool {
	return ifi.Flags&net.FlagMulticast != 0
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
