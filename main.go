package main

import (
	"crypto/tls"
	"flag"
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
	timeout := flag.Int("timeout", 5, "connect timeout")
	flag.Parse()

	log.SetFlags(log.Lshortfile | log.LstdFlags)

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
		Timeout: time.Duration(*timeout) * time.Second,
	}

	// Display connection parameters
	log.Printf("Connecting to:   %s\n", addr)
	log.Printf("TLS server name: %s\n", tlsConfig.ServerName)
	log.Printf("Insecure TLS:    %v\n", *insecure)

	// Establish TLS connection
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	if err := conn.Handshake(); err != nil {
		log.Printf("TLS handshake failed: %v", err)
		return
	}
	log.Printf("Connection Handshake succeeded\n")

	// Retrieve and print peer certificates
	state := conn.ConnectionState()
	for i, cert := range state.PeerCertificates {
		log.Printf("Certificate #%d: Subject: %s\n", i+1, cert.Subject)
		log.Printf("  Subject: %s\n", cert.Subject.String())
		log.Printf("  Issuer:  %s\n", cert.Issuer.String())
		log.Printf("  NotBefore: %s\n", cert.NotBefore)
		log.Printf("  NotAfter:  %s\n", cert.NotAfter)
	}

	// TLS handshake info
	log.Printf("TLS handshake complete: %v\n", state.HandshakeComplete)
	log.Printf("Mutual protocol negotiated: %v\n", state.NegotiatedProtocolIsMutual)
}
