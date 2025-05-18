package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	// CLI flags
	host := flag.String("host", "localhost", "Target hostname or IP address")
	port := flag.String("port", "443", "Target TCP port (default 443)")
	serverName := flag.String("servername", "", "SNI server name (defaults to host)")
	insecure := flag.Bool("insecure", false, "Skip TLS certificate verification (not recommended)")
	flag.Parse()

	log.SetFlags(log.Lshortfile)

	// Compose address
	addr := net.JoinHostPort(*host, *port)

	// Determine server name (SNI)
	sni := *serverName
	if sni == "" {
		sni = *host
	}

	// Create a TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: *insecure,
		ServerName:         sni,
	}

	// Use modern Dialer (Go 1.24+; DualStack removed)
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	// Display connection parameters
	fmt.Printf("Connecting to:   %s\n", addr)
	fmt.Printf("TLS server name: %s\n", tlsConfig.ServerName)
	fmt.Printf("Insecure TLS:    %v\n", *insecure)

	// Establish TLS connection
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Retrieve and print peer certificates
	state := conn.ConnectionState()
	for i, cert := range state.PeerCertificates {
		fmt.Printf("Certificate %d Subject: %s\n", i+1, cert.Subject)
	}

	// TLS handshake info
	log.Printf("TLS handshake complete: %v", state.HandshakeComplete)
	log.Printf("Mutual protocol negotiated: %v", state.NegotiatedProtocolIsMutual)
}
