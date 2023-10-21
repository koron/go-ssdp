package main

import (
	"flag"
	"log"
	"os"

	"github.com/koron/go-ssdp"
)

func main() {
	nt := flag.String("nt", "my:device", "NT: Type")
	usn := flag.String("usn", "unique:id", "USN: ID")
	loc := flag.String("loc", "", "LOCATION: location header")
	srv := flag.String("srv", "", "SERVER:  server header")
	maxAge := flag.Int("maxage", 1800, "cache control, max-age")
	laddr := flag.String("laddr", "", "local address to listen")
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
	if *ttl > 0 {
		ssdp.SetMulticastTTL(*ttl)
	}
	if *sysIf {
		ssdp.SetMulticastSystemAssignedInterface(true)
	}

	err := ssdp.AnnounceAlive(*nt, *usn, *loc, *srv, *maxAge, *laddr)
	if err != nil {
		log.Fatal(err)
	}
}
