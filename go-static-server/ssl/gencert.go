// Package main generates a self-signed TLS certificate and private key.
// It is compiled and run during the Docker build to produce certs
// without requiring openssl in the final image.
//
// Usage: gencert [-cert path] [-key path] [-duration days]
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	certPath := flag.String("cert", "/etc/ssl/server.crt", "output path for the certificate PEM")
	keyPath := flag.String("key", "/etc/ssl/server.key", "output path for the private key PEM")
	days := flag.Int("duration", 3650, "certificate validity in days")
	flag.Parse()

	// Generate ECDSA P-256 private key.
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate private key: %v\n", err)
		os.Exit(1)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate serial number: %v\n", err)
		os.Exit(1)
	}

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"cloud-bulldozer"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(time.Duration(*days) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create certificate: %v\n", err)
		os.Exit(1)
	}

	// Write certificate PEM.
	certFile, err := os.Create(*certPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create cert file: %v\n", err)
		os.Exit(1)
	}
	defer certFile.Close()
	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write cert PEM: %v\n", err)
		os.Exit(1)
	}

	// Write private key PEM.
	keyFile, err := os.Create(*keyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create key file: %v\n", err)
		os.Exit(1)
	}
	defer keyFile.Close()
	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal private key: %v\n", err)
		os.Exit(1)
	}
	if err := pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write key PEM: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Certificate written to %s\n", *certPath)
	fmt.Printf("Private key written to %s\n", *keyPath)
}
