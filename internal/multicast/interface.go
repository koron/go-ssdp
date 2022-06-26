package multicast

import (
	"net"
	"sync"
)

type InterfacesProviderFunc func() []net.Interface

// InterfacesProvider specify a function to list all interfaces to multicast.
// If no provider are given, all possible interfaces will be used.
var InterfacesProvider InterfacesProviderFunc

var ifLock sync.Mutex
var ifList []net.Interface

// interfaces gets list of net.Interface to multicast UDP packet.
func interfaces() ([]net.Interface, error) {
	ifLock.Lock()
	defer ifLock.Unlock()
	if InterfacesProvider != nil {
		if list := InterfacesProvider(); len(list)>0 {
			return list, nil
		}
	}
	if len(ifList) > 0 {
		return ifList, nil
	}
	l, err := interfacesIPv4()
	if err != nil {
		return nil, err
	}
	ifList = l
	return ifList, nil
}

// interfacesIPv4 lists net.Interface on IPv4.
func interfacesIPv4() ([]net.Interface, error) {
	iflist, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	list := make([]net.Interface, 0, len(iflist))
	for _, ifi := range iflist {
		if !hasLinkUp(&ifi) || !hasIPv4Address(&ifi) {
			continue
		}
		list = append(list, ifi)
	}
	return list, nil
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
