package main

import (
	"crypto/tls"
	//"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	host := flag.String("host", "localhost", "hostname")
	port := flag.String("port", "443", "port")
	servername := flag.String("servername", "", "servername")
	insecure := flag.Bool("insecure", false, "insecure or not")
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	dialer := &net.Dialer{
		Timeout:   5000 * time.Millisecond,
		DualStack: false,
		LocalAddr: nil,
	}
	addr := *host + ":" + *port

	myservername := *servername
	if myservername == "" {
		myservername = *host
	}

	conf := &tls.Config{
		InsecureSkipVerify: *insecure,
		ServerName:         myservername,
	}

	fmt.Printf("addrress: [%v]\n", addr)
	fmt.Printf("servername: [%v]\n", conf.ServerName)
	fmt.Printf("insecure [%v]\n", *insecure)

	conn, err := tls.DialWithDialer(dialer, "tcp4", addr, conf)

	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	defer conn.Close()
	if err != nil && strings.Contains(err.Error(), "timed out") {
		fmt.Println("Connection timed out")
		os.Exit(3)
	}

	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		//fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Println(v.Subject)
	}
	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

	conn.Close()
}
