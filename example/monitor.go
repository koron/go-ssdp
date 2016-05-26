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

	_, err := ssdp.NewMonitor(aliveHandler, byeHandler)
	if err != nil {
		log.Fatal(err)
	}
	for {}
}

func aliveHandler(m *ssdp.Alive) {
	log.Printf("Alive: %#v", m)
}

func byeHandler(m *ssdp.Bye) {
	log.Printf("Bye: %#v", m)
}
