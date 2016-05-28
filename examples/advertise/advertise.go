package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/koron/go-ssdp"
)

func main() {
	st := flag.String("st", "my:device", "ST: Type")
	usn := flag.String("usn", "unique:id", "USN: ID")
	loc := flag.String("loc", "", "LOCATION: location header")
	srv := flag.String("srv", "", "SERVER:  server header")
	maxAge := flag.Int("maxage", 1800, "cache control, max-age")
	ai := flag.Int("ai", 0, "alive interval")
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

	ad, err := ssdp.Advertise(*st, *usn, *loc, *srv, *maxAge)
	if err != nil {
		log.Fatal(err)
	}
	var (
		quit = make(chan bool)
		ac <-chan time.Time
	)
	if *ai > 0 {
		ac = time.Tick(time.Duration(*ai) * time.Second)
	} else {
		ac = make(chan time.Time)
	}

loop:
	for {
		select {
		case <-ac:
			ad.Alive()
		case <-quit:
			break loop
		}
	}
}
