package ssdp

import (
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAdvertise_Alive(t *testing.T) {
	var mu sync.Mutex
	var mm []*AliveMessage
	m := newTestMonitor("test:advertise+alive", func(m *AliveMessage) {
		mu.Lock()
		mm = append(mm, m)
		mu.Unlock()
	}, nil, nil)
	err := m.Start()
	if err != nil {
		t.Fatalf("failed to start Monitor: %s", err)
	}
	defer m.Close()

	a, err := Advertise("test:advertise+alive", "usn:advertise+alive", "location:advertise+alive", "server:advertise+alive", 600)
	if err != nil {
		t.Fatalf("failed to Advertise: %s", err)
	}
	err = a.Alive()
	if err != nil {
		a.Close()
		t.Fatalf("failed to send alive: %s", err)
	}
	a.Close()
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

	expHdr := map[string]string{
		"Nts":           "ssdp:alive",
		"Nt":            "test:advertise+alive",
		"Usn":           "usn:advertise+alive",
		"Location":      "location:advertise+alive",
		"Server":        "server:advertise+alive",
		"Cache-Control": "max-age=600",
	}
	for i, m := range mm {
		if strings.HasSuffix(port, m.From.String()) {
			t.Errorf("unmatch port#%d:\nwant=%q\n got=%q", i, port, m.From.String())
		}
		if m.Type != "test:advertise+alive" {
			t.Errorf("unexpected alive#%d type: want=%q got=%q", i, "test:advertise+alive", m.Type)
		}
		if m.USN != "usn:advertise+alive" {
			t.Errorf("unexpected alive#%d usn: want=%q got=%q", i, "usn:advertise+alive", m.USN)
		}
		if m.Location != "location:advertise+alive" {
			t.Errorf("unexpected alive#%d location: want=%q got=%q", i, "location:advertise+alive", m.Location)
		}
		if m.Server != "server:advertise+alive" {
			t.Errorf("unexpected alive#%d server: want=%q got=%q", i, "server:advertise+alive", m.Server)
		}
		if m.MaxAge() != 600 {
			t.Errorf("unexpected max-age: want=%d got=%d", 600, m.MaxAge())
		}

		h := m.Header()
		for k := range h {
			exp, ok := expHdr[k]
			if !ok {
				t.Errorf("unexpected header #%d %q=%q", i, k, h.Get(k))
			} else if act := h.Get(k); act != exp {
				t.Errorf("header #%d %q value mismatch:\nwant=%q\n got=%q", i, k, exp, act)
			}
		}
	}
}

func TestAdvertise_Bye(t *testing.T) {
	var mu sync.Mutex
	var mm []*ByeMessage
	m := newTestMonitor("test:advertise+bye", nil, func(m *ByeMessage) {
		mu.Lock()
		mm = append(mm, m)
		mu.Unlock()
	}, nil)
	err := m.Start()
	if err != nil {
		t.Fatalf("failed to start Monitor: %s", err)
	}
	defer m.Close()

	a, err := Advertise("test:advertise+bye", "usn:advertise+bye", "location:advertise+bye", "server:advertise+bye", 600)
	if err != nil {
		t.Fatalf("failed to Advertise: %s", err)
	}
	err = a.Bye()
	if err != nil {
		a.Close()
		t.Fatalf("failed to send bye: %s", err)
	}
	a.Close()
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

	expHdr := map[string]string{
		"Nts": "ssdp:byebye",
		"Nt":  "test:advertise+bye",
		"Usn": "usn:advertise+bye",
	}
	for i, m := range mm {
		if strings.HasSuffix(port, m.From.String()) {
			t.Errorf("unmatch port#%d:\nwant=%q\n got=%q", i, port, m.From.String())
		}
		if m.Type != "test:advertise+bye" {
			t.Errorf("unexpected bye#%d type: want=%q got=%q", i, "test:advertise+bye", m.Type)
		}
		if m.USN != "usn:advertise+bye" {
			t.Errorf("unexpected bye#%d usn: want=%q got=%q", i, "usn:advertise+bye", m.USN)
		}

		h := m.Header()
		for k := range h {
			exp, ok := expHdr[k]
			if !ok {
				t.Errorf("unexpected header #%d %q=%q", i, k, h.Get(k))
			} else if act := h.Get(k); act != exp {
				t.Errorf("header #%d %q value mismatch:\nwant=%q\n got=%q", i, k, exp, act)
			}
		}
	}
}
