package ssdp

import "testing"

func TestSetMulticastRecvAddrIPv4(t *testing.T) {
	_, err := recvAddrResolver.resolve()
	if err != nil {
		t.Errorf("resolve #1 failed: %s", err)
	}
	SetMulticastRecvAddrIPv4("224.0.0.0:1900")
	if recvAddrResolver.addr != "224.0.0.0:1900" {
		t.Errorf("unexpected recvAddrResolver.addr:\nwant=%q got=%q", "224.0.0.0:1900", recvAddrResolver.addr)
	}
	_, err = recvAddrResolver.resolve()
	if err != nil {
		t.Errorf("resolve #2 failed: %s", err)
	}
}

func TestSetMulticastSendAddrIPv4(t *testing.T) {
	_, err := sendAddrResolver.resolve()
	if err != nil {
		t.Errorf("resolve #1 failed: %s", err)
	}
	SetMulticastSendAddrIPv4("239.255.255.250:1900")
	if sendAddrResolver.addr != "239.255.255.250:1900" {
		t.Errorf("unexpected sendAddrResolver.addr:\nwant=%q got=%q", "239.255.255.250:1900", sendAddrResolver.addr)
	}
	_, err = sendAddrResolver.resolve()
	if err != nil {
		t.Errorf("resolve #2 failed: %s", err)
	}
}
