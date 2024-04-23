package ssdp

import (
	"sync"
	"testing"
	"time"
)

// TestTTL tests TTL() doesn't something bad.
func TestTTL(t *testing.T) {
	// start alive monitor.
	var mu sync.Mutex
	var alives []*AliveMessage
	m := newTestMonitor("test:ttl", func(m *AliveMessage) {
		mu.Lock()
		alives = append(alives, m)
		mu.Unlock()
	}, nil, nil)
	err := m.Start()
	if err != nil {
		t.Fatalf("failed to start Monitor: %s", err)
	}
	defer m.Close()

	// send test alive with TTL:2
	err = AnnounceAlive("test:ttl", "usn:ttl", "location:ttl", "server:ttl", 600, "", TTL(2))
	if err != nil {
		t.Fatalf("failed to announce alive: %s", err)
	}
	time.Sleep(500 * time.Millisecond)

	checkAlives(t, alives, "test:ttl", "usn:ttl", "location:ttl", "server:ttl")
}

// TestOnlySystemInterface tests OnlySystemInterface().
// Monitor with OnlySystemInterface() and send alive message to all interfaces.
// Monitor will receive just an alive message for default interface.
func TestOnlySystemInterface(t *testing.T) {
	// start alive monitor with OnlySystemInterface.
	var mu sync.Mutex
	var alives []*AliveMessage
	m := newTestMonitor("test:sysif", func(m *AliveMessage) {
		mu.Lock()
		alives = append(alives, m)
		mu.Unlock()
	}, nil, nil)
	m.Options = append(m.Options, OnlySystemInterface())
	err := m.Start()
	if err != nil {
		t.Fatalf("failed to start Monitor: %s", err)
	}
	defer m.Close()

	// send a test alive.
	err = AnnounceAlive("test:sysif", "usn:sysif", "location:sysif", "server:sysif", 600, "")
	if err != nil {
		t.Fatalf("failed to announce alive: %s", err)
	}
	time.Sleep(500 * time.Millisecond)

	checkAlives(t, alives, "test:sysif", "usn:sysif", "location:sysif", "server:sysif")

	if len(alives) != 1 {
		t.Fatalf("exact an alive should be detected: but got %d", len(alives))
	}
}

func TestLocalAddr(t *testing.T) {
	// TODO:
}

func TestRemoteAddr(t *testing.T) {
	// TODO:
}
