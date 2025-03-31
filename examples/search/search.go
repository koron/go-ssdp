package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/koron/go-ssdp"
)

func main() {
	t := flag.String("t", ssdp.All, "search type")
	w := flag.Uint("w", 1, "wait time")
	l := flag.String("l", "", "local address to listen")
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

	list, err := ssdp.Search(*t, *w, *l, opts...)
	if err != nil {
		log.Fatal(err)
	}
	for i, srv := range list {
		//fmt.Printf("%d: %#v\n", i, srv)
		fmt.Printf("%d: %s %s\n", i, srv.Type, srv.Location)
	}
}
