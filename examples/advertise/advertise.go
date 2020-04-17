package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/koron/go-ssdp"
)

func main() {
	st := flag.String("st", "my:device", "ST: Type")
	usn := flag.String("usn", "unique:id", "USN: ID")
	loc := flag.String("loc", "", "LOCATION: location header")
	srv := flag.String("srv", "", "SERVER:  server header")
	maxAge := flag.Int("maxage", 1800, "cache control, max-age")
	ai := flag.Int("ai", 10, "alive interval")
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

	locationHandler := func(originator net.Addr) string {
		if originator != nil {
			ssdp.Logger.Printf("Location header handler called in the context of an M-SEARCH from %s", originator.String())
			// do something with the originator argument to compute the location
		} else {
			ssdp.Logger.Printf("Location header handler called in the context of an Alive announcement")
		}
		// just returning the location argument
		return *loc
	}

	ad, err := ssdp.Advertise(*st, *usn, *srv, *maxAge, locationHandler)
	if err != nil {
		log.Fatal(err)
	}
	var aliveTick <-chan time.Time
	if *ai > 0 {
		aliveTick = time.Tick(time.Duration(*ai) * time.Second)
	} else {
		aliveTick = make(chan time.Time)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

loop:
	for {
		select {
		case <-aliveTick:
			ad.Alive()
		case <-quit:
			break loop
		}
	}
	ad.Bye()
	ad.Close()
}
