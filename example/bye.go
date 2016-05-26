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
	laddr := flag.String("laddr", "", "local address to listen")
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

	err := ssdp.AnnounceBye(*nt, *usn, *laddr)
	if err != nil {
		log.Fatal(err)
	}
}
