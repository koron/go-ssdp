package ssdp

import (
	"fmt"
	"net"
	"testing"
)

func TestMaxAgeDirectiveCase(t *testing.T) {
	for _, test := range []struct {
		value string
		want  int
	}{
		{"max-age=1800", 1800},
		{"MAX-AGE=1800", 1800},
		{"Max-Age = 1800", 1800},
		{"public, mAx-AgE=1800", 1800},
		{"MAX-AGE=0", 0},
		{"MAX-AGE=-1", -1},
		{"MAX-AGE=2147483648", -1},
		{"xmax-age=1800", -1},
		{"no-cache", -1},
		{"", -1},
	} {
		t.Run(test.value, func(t *testing.T) {
			response := fmt.Sprintf("HTTP/1.1 200 OK\r\nST: test-device\r\nUSN: test-id\r\nCACHE-CONTROL: %s\r\nContent-Length: 0\r\n\r\n", test.value)
			service, err := parseService([]byte(response))
			if err != nil {
				t.Fatal(err)
			}
			if service.Type != "test-device" || service.USN != "test-id" {
				t.Fatalf("unexpected service: %+v", service)
			}
			if got := service.MaxAge(); got != test.want {
				t.Errorf("service MaxAge = %d; want %d", got, test.want)
			}
			var message *AliveMessage
			monitor := &Monitor{Alive: func(m *AliveMessage) { message = m }}
			notify := fmt.Sprintf("NOTIFY * HTTP/1.1\r\nNT: test-device\r\nNTS: ssdp:alive\r\nUSN: test-id\r\nCACHE-CONTROL: %s\r\n\r\n", test.value)
			if err := monitor.handleRaw(&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1900}, []byte(notify)); err != nil {
				t.Fatal(err)
			}
			if message == nil || message.Type != "test-device" || message.USN != "test-id" {
				t.Fatalf("unexpected alive notification: %+v", message)
			}
			if got := message.MaxAge(); got != test.want {
				t.Errorf("alive MaxAge = %d; want %d", got, test.want)
			}
		})
	}
}
