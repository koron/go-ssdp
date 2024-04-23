package ssdp

import (
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/exp/slices"
)

func checkAlives(t *testing.T, alives []*AliveMessage, typ, usn, loc, srv string) {
	t.Helper()

	first := slices.IndexFunc(alives, func(m *AliveMessage) bool {
		return m.Type == typ
	})
	if first < 0 {
		t.Errorf("no AliveMessage which Type is %s", typ)
		return
	}

	_, port, err := net.SplitHostPort(alives[first].From.String())
	if err != nil {
		t.Errorf("failed to split host and port from first message: %s", err)
		return
	}
	port = ":" + port

	for i, m := range alives {
		if m.Type != typ {
			t.Logf("unexpected alive[%d].Type: want=%q got=%q", i, typ, m.Type)
			continue
		}
		if !strings.HasSuffix(m.From.String(), port) {
			t.Errorf("unmatch alive[%d].From (:port): want=%q got=%q", i, port, m.From.String())
		}
		if m.USN != usn {
			t.Errorf("unexpected alive[%d].USN: want=%q got=%q", i, usn, m.USN)
		}
		if m.Location != loc {
			t.Errorf("unexpected alive[%d].Location: want=%q got=%q", i, loc, m.Location)
		}
		if m.Server != srv {
			t.Errorf("unexpected alive[%d].Server: want=%q got=%q", i, srv, m.Server)
		}
	}
}

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

	checkAlives(t, mm, "test:announce+alive", "usn:announce+alive", "location:announce+alive", "server:announce+alive")
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
