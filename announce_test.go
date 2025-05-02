package ssdp

import (
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAnnounceAlive(t *testing.T) {
	var mu sync.Mutex
	var mm []*AliveMessage
	m := newTestMonitor("test:announce+alive", func(m *AliveMessage) {
		mu.Lock()
		mm = append(mm, m)
		mu.Unlock()
	}, nil, nil)
	err := m.Start()
	if err != nil {
		t.Fatalf("failed to start Monitor: %s", err)
	}
	defer m.Close()

	err = AnnounceAlive("test:announce+alive", "usn:announce+alive", "location:announce+alive", "server:announce+alive", 600, "")
	if err != nil {
		t.Fatalf("failed to announce alive: %s", err)
	}
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(mm) < 1 {
		t.Fatal("no alives detected")
	}
	//t.Logf("found %d alives", len(mm))
	_, port, err := net.SplitHostPort(mm[0].From.String())
	if err != nil {
		t.Fatalf("failed to split host and port: %s", err)
	}
	port = ":" + port
	for i, m := range mm {
		if strings.HasSuffix(port, m.From.String()) {
			t.Errorf("unmatch port#%d:\nwant=%q\n got=%q", i, port, m.From.String())
		}
		if m.Type != "test:announce+alive" {
			t.Errorf("unexpected alive#%d type: want=%q got=%q", i, "test:announce+alive", m.Type)
		}
		if m.USN != "usn:announce+alive" {
			t.Errorf("unexpected alive#%d usn: want=%q got=%q", i, "usn:announce+alive", m.USN)
		}
		if m.Location != "location:announce+alive" {
			t.Errorf("unexpected alive#%d location: want=%q got=%q", i, "location:announce+alive", m.Location)
		}
		if m.Server != "server:announce+alive" {
			t.Errorf("unexpected alive#%d server: want=%q got=%q", i, "server:announce+alive", m.Server)
		}
	}
}

func TestAnnounceBye(t *testing.T) {
	var mu sync.Mutex
	var mm []*ByeMessage
	m := newTestMonitor("test:announce+bye", nil, func(m *ByeMessage) {
		mu.Lock()
		mm = append(mm, m)
		mu.Unlock()
	}, nil)
	err := m.Start()
	if err != nil {
		t.Fatalf("failed to start Monitor: %s", err)
	}
	defer m.Close()

	err = AnnounceBye("test:announce+bye", "usn:announce+bye", "")
	if err != nil {
		t.Fatalf("failed to announce bye: %s", err)
	}
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(mm) < 1 {
		t.Fatal("no byes detected")
	}
	//t.Logf("found %d byes", len(mm))
	_, port, err := net.SplitHostPort(mm[0].From.String())
	if err != nil {
		t.Fatalf("failed to split host and port: %s", err)
	}
	port = ":" + port
	for i, m := range mm {
		if strings.HasSuffix(port, m.From.String()) {
			t.Errorf("unmatch port#%d:\nwant=%q\n got=%q", i, port, m.From.String())
		}
		if m.Type != "test:announce+bye" {
			t.Errorf("unexpected bye#%d type: want=%q got=%q", i, "test:announce+bye", m.Type)
		}
		if m.USN != "usn:announce+bye" {
			t.Errorf("unexpected bye#%d usn: want=%q got=%q", i, "usn:announce+bye", m.USN)
		}
	}
}
