package multicast

import "testing"

func TestSetMulticastRecvAddrIPv4(t *testing.T) {
	_, err := RecvAddrResolver.resolve()
	if err != nil {
		t.Errorf("resolve #1 failed: %s", err)
	}
	SetRecvAddrIPv4("224.0.0.0:1900")
	if RecvAddrResolver.Addr != "224.0.0.0:1900" {
		t.Errorf("unexpected RecvAddrResolver.Addr:\nwant=%q got=%q", "224.0.0.0:1900", RecvAddrResolver.Addr)
	}
	_, err = RecvAddrResolver.resolve()
	if err != nil {
		t.Errorf("resolve #2 failed: %s", err)
	}
}

func TestSetMulticastSendAddrIPv4(t *testing.T) {
	_, err := sendAddrResolver.resolve()
	if err != nil {
		t.Errorf("resolve #1 failed: %s", err)
	}
	SetSendAddrIPv4("239.255.255.250:1900")
	if sendAddrResolver.Addr != "239.255.255.250:1900" {
		t.Errorf("unexpected sendAddrResolver.Addr:\nwant=%q got=%q", "239.255.255.250:1900", sendAddrResolver.Addr)
	}
	_, err = sendAddrResolver.resolve()
	if err != nil {
		t.Errorf("resolve #2 failed: %s", err)
	}
}
