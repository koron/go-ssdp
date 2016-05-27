package main

import (
	"flag"
	"log"
	"os"

	"github.com/koron/go-ssdp"
)

func main() {
	v := flag.Bool("v", false, "verbose mode")
	h := flag.Bool("h", false, "show help")
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	if *v {
		ssdp.Logger = log.New(os.Stderr, "[SSDP] ", log.LstdFlags)
	}

	_, err := ssdp.NewMonitor(aliveHandler, byeHandler, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
	}
}

func aliveHandler(m *ssdp.Alive) {
	log.Printf("Alive: From=%s Type=%s USN=%s Location=%s Server=%s MaxAge=%d",
		m.From.String(), m.Type, m.USN, m.Location, m.Server, m.MaxAge())
}

func byeHandler(m *ssdp.Bye) {
	log.Printf("Bye: From=%s Type=%s USN=%s", m.From.String(), m.Type, m.USN)
}
