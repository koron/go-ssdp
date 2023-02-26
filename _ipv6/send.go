// IPv6 マルチキャストを送信する
//
// Windowsでも動いた
// アドレスに %{インターフェス名} を含めることでnet.Interfaceを指定してる
package main

import (
	"log"
	"net"
)

func main() {
	if err := send(); err != nil {
		log.Fatal(err)
	}
}

//const addrStr = "239.255.255.250:1900"
const addrStr = "[FF02::C%イーサネット 2]:1900"

func send() error {
	//addr, err := net.ResolveUDPAddr("udp", "")
	//if err != nil {
	//	return err
	//}
	//log.Printf("addr=%s", addr)
	conn, err := net.Dial("udp", addrStr)
	if err != nil {
		return err
	}
	defer conn.Close()
	n, err := conn.Write([]byte("NOTIFY "))
	if err != nil {
		return err
	}
	log.Printf("sent %d bytes", n)
	return nil
}
