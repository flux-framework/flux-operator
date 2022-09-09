package certs

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

// NewCertificate generates a new certificate
func NewCertificate(hosts []string, isCa bool) *Certificate {
	return &Certificate{Hosts: hosts, IsCA: isCa}
}

type Certificate struct {
	Hosts []string

	// see https://www.rfc-editor.org/rfc/rfc5280#section-4.2.1.9
	IsCA    bool
	Public  string
	Private string
}

// pemBlockForKey generates the private key block
func pemBlockForKey(private *ecdsa.PrivateKey) *pem.Block {

	b, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
		os.Exit(2)
	}
	return &pem.Block{Type: "PRIVATE KEY", Bytes: b}
}

// Generate a serial number
func getSerialNumber() *big.Int {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	number, err := rand.Int(rand.Reader, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue generating serial number for certificate: %v", err)
		os.Exit(2)
	}
	return number
}

func publicKey(priv any) any {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil

	}
}

// Generate makes a certificate (e.g., key and secret)
// If they are ephemeral, I think we can do this
func (c *Certificate) Generate() {

	// https://go.dev/src/crypto/tls/generate_cert.go
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	template := x509.Certificate{
		SerialNumber: getSerialNumber(),
		Subject: pkix.Name{
			Organization: []string{"Flux Framework"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 30),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, host := range c.Hosts {

		// Does the host have an ip address or a DNS name?
		if ip := net.ParseIP(host); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, host)
		}

	}
	// Can this certificate be used to verify other certificate signatures?
	// Otherwise, it's a leaf (false)
	if c.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(private), private)
	if err != nil {
		log.Fatalf("Failed to create public certificate: %s", err)
	}
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	c.Public = out.String()
	out.Reset()

	// Generate private certificate
	pem.Encode(out, pemBlockForKey(private))
	c.Private = out.String()
}
