package main

import (
	"flag"
	"log"
	"os"

	"github.com/koron/go-ssdp"
)

func main() {
	st := flag.String("st", "my:device", "ST: Type")
	usn := flag.String("usn", "unique:id", "USN: ID")
	loc := flag.String("loc", "", "LOCATION: location header")
	srv := flag.String("srv", "", "SERVER:  server header")
	maxAge := flag.Int("maxage", 1800, "cache control, max-age")
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

	_, err := ssdp.Advertise(*st, *usn, *loc, *srv, *maxAge)
	if err != nil {
		log.Fatal(err)
	}
	for {
	}
}
