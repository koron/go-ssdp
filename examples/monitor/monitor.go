package main

import (
	"flag"
	"log"
	"os"

	"github.com/koron/go-ssdp"
)

var st string

func main() {
	v := flag.Bool("v", false, "verbose mode")
	h := flag.Bool("h", false, "show help")
	flag.StringVar(&st, "st", "", "ST: Type")

	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	if *v {
		ssdp.Logger = log.New(os.Stderr, "[SSDP] ", log.LstdFlags)
	}

	m := &ssdp.Monitor{
		Alive:  onAlive,
		Bye:    onBye,
		Search: onSearch,
	}
	if err := m.Start(); err != nil {
		log.Fatal(err)
	}
	// wait infinitely
	ch := make(chan struct{})
	<-ch
}

func onAlive(m *ssdp.AliveMessage) {
	if st != "" && st != m.Type { return }

	log.Printf("Alive: From=%s Type=%s USN=%s Location=%s Server=%s MaxAge=%d",
		m.From.String(), m.Type, m.USN, m.Location, m.Server, m.MaxAge())
}

func onBye(m *ssdp.ByeMessage) {
	if st != "" && st != m.Type { return }

	log.Printf("Bye: From=%s Type=%s USN=%s", m.From.String(), m.Type, m.USN)
}

func onSearch(m *ssdp.SearchMessage) {
	if st != "" && st != m.Type { return }

	log.Printf("Search: From=%s Type=%s", m.From.String(), m.Type)
}
