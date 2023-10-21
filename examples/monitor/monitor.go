package main

import (
	"flag"
	"log"
	"os"

	"github.com/koron/go-ssdp"
)

var filterType string

func main() {
	v := flag.Bool("v", false, "verbose mode")
	h := flag.Bool("h", false, "show help")
	flag.StringVar(&filterType, "filter_type", "", "print only a specified type (ST or NT). default is print all types.")
	ttl := flag.Int("ttl", 0, "TTL for outgoing multicast packets")
	sysIf := flag.Bool("sysif", false, "use system assigned multicast interface")
	flag.Parse()

	if *h {
		flag.Usage()
		return
	}
	if *v {
		ssdp.Logger = log.New(os.Stderr, "[SSDP] ", log.LstdFlags)
	}

	var opts []ssdp.Option
	if *ttl > 0 {
		opts = append(opts, ssdp.TTL(*ttl))
	}
	if *sysIf {
		opts = append(opts, ssdp.OnlySystemInterface())
	}

	m := &ssdp.Monitor{
		Alive:   onAlive,
		Bye:     onBye,
		Search:  onSearch,
		Options: opts,
	}
	if err := m.Start(); err != nil {
		log.Fatal(err)
	}
	// wait infinitely
	ch := make(chan struct{})
	<-ch
}

// filterByType returns true when given type must be hidden.
func filterByType(typ string) bool {
	if filterType != "" && filterType != typ {
		return true
	}
	return false
}

func onAlive(m *ssdp.AliveMessage) {
	if filterByType(m.Type) {
		return
	}

	log.Printf("Alive: From=%s Type=%s USN=%s Location=%s Server=%s MaxAge=%d",
		m.From.String(), m.Type, m.USN, m.Location, m.Server, m.MaxAge())
}

func onBye(m *ssdp.ByeMessage) {
	if filterByType(m.Type) {
		return
	}

	log.Printf("Bye: From=%s Type=%s USN=%s", m.From.String(), m.Type, m.USN)
}

func onSearch(m *ssdp.SearchMessage) {
	if filterByType(m.Type) {
		return
	}

	log.Printf("Search: From=%s Type=%s", m.From.String(), m.Type)
}
