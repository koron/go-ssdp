package main

import (
	"flag"
	"log"
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
	ttl := flag.Int("ttl", 0, "TTL for outgoing multicast packets")
	sysIf := flag.Bool("sysif", false, "use system assigned multicast interface")
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

	var opts []ssdp.Option
	if *ttl > 0 {
		opts = append(opts, ssdp.TTL(*ttl))
	}
	if *sysIf {
		opts = append(opts, ssdp.OnlySystemInterface())
	}

	ad, err := ssdp.Advertise(*st, *usn, *loc, *srv, *maxAge, opts...)
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
