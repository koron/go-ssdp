package multicast

import "testing"

func TestSetMulticastRecvAddrIPv4(t *testing.T) {
	// resolve with default
	_, err := RecvAddrResolver.resolve()
	if err != nil {
		t.Errorf("resolve #1 failed: %s", err)
	}

	// resolve after override
	SetRecvAddrIPv4("224.0.0.1:1900")
	if RecvAddrResolver.Addr != "224.0.0.1:1900" {
		t.Errorf("unexpected RecvAddrResolver.Addr:\nwant=%q got=%q", "224.0.0.1:1900", RecvAddrResolver.Addr)
	}
	_, err = RecvAddrResolver.resolve()
	if err != nil {
		t.Errorf("resolve #2 failed: %s", err)
	}
}

func TestSetMulticastSendAddrIPv4(t *testing.T) {
	// resolve with default
	_, err := SendAddr()
	if err != nil {
		t.Errorf("resolve #1 failed: %s", err)
	}

	// resolve after override
	SetSendAddrIPv4("239.255.255.250:1900")
	if sendAddrResolver.Addr != "239.255.255.250:1900" {
		t.Errorf("unexpected sendAddrResolver.Addr:\nwant=%q got=%q", "239.255.255.250:1900", sendAddrResolver.Addr)
	}
	first, err := SendAddr()
	if err != nil {
		t.Errorf("resolve #2 failed: %s", err)
	}

	// resolve by cache
	second, err := SendAddr()
	if err != nil {
		t.Errorf("resolve #3 failed: %s", err)
	}
	if second != first {
		t.Errorf("cache mismatch: first=%p second=%p", first, second)
	}
}
