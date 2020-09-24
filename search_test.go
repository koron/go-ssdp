package ssdp

import (
	"net"
	"strings"
	"sync"
	"testing"
)

func testMaxAge(t *testing.T, s string, expect int) {
	act := extractMaxAge(s, -1)
	if act != expect {
		t.Errorf("max-age for %q should be %d but actually %d", s, expect, act)
	}
}

func TestExtractMaxAge(t *testing.T) {
	// empty
	testMaxAge(t, "", -1)
	// spaces around `=`
	testMaxAge(t, "max-age=100", 100)
	testMaxAge(t, "max-age = 200", 200)
	testMaxAge(t, "max-age= 300", 300)
	testMaxAge(t, "max-age =400", 400)
	// minus
	testMaxAge(t, "max-age=-100", -1)
	// invalid name
	testMaxAge(t, "foo=100", -1)
	// contained valid name
	testMaxAge(t, "foomax-age=100", -1)
	// surrounded
	testMaxAge(t, ";max-age=500;", 500)
	testMaxAge(t, ";max-age=600", 600)
	testMaxAge(t, "max-age=700;", 700)
}

func TestSearch_Request(t *testing.T) {
	searchType := "test:search+request"

	var mu sync.Mutex
	var mm []*SearchMessage
	m := newTestMonitor(searchType, nil, nil, func(m *SearchMessage) {
		mu.Lock()
		mm = append(mm, m)
		mu.Unlock()
	})
	err := m.Start()
	if err != nil {
		t.Fatalf("failed to start Monitor: %s", err)
	}
	defer m.Close()

	srvs, err := Search(searchType, 1, "")
	if err != nil {
		t.Fatalf("failed to Search: %s", err)
	}
	if len(srvs) > 0 {
		t.Errorf("unexpected services: %+v", srvs)
	}

	if len(mm) < 1 {
		t.Fatal("no search detected")
	}
	_, port, err := net.SplitHostPort(mm[0].From.String())
	if err != nil {
		t.Fatalf("failed to split host and port: %s", err)
	}
	port = ":" + port
	for i, m := range mm {
		if m.Type != "test:search+request" {
			t.Errorf("unmatch type#%d:\nwant=%q\n got=%q", i, "test:search+request", m.Type)
		}
		if strings.HasSuffix(port, m.From.String()) {
			t.Errorf("unmatch port#%d:\nwant=%q\n got=%q", i, port, m.From.String())
		}
	}
}

func TestSearch_Response(t *testing.T) {
	a, err := Advertise("test:search+response", "usn:search+response", "location:search+response", "server:search+response", 600)
	if err != nil {
		t.Fatalf("failed to Advertise: %s", err)
	}
	defer a.Close()

	srvs, err := Search("test:search+response", 1, "")
	if err != nil {
		t.Fatalf("failed to Search: %s", err)
	}
	if len(srvs) == 0 {
		t.Errorf("no services found")
	}

	for i, s := range srvs {
		if s.Type != "test:search+response" {
			t.Errorf("unexpected service#%d type: want=%q got=%q", i, "test:search+response", s.Type)
		}
		if s.USN != "usn:search+response" {
			t.Errorf("unexpected service#%d usn: want=%q got=%q", i, "usn:search+response", s.USN)
		}
		if s.Location != "location:search+response" {
			t.Errorf("unexpected service#%d location: want=%q got=%q", i, "location:search+response", s.Location)
		}
		if s.Server != "server:search+response" {
			t.Errorf("unexpected service#%d server: want=%q got=%q", i, "server:search+response", s.Server)
		}
	}
}
