package certs

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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
	b, err := x509.MarshalECPrivateKey(private)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
		os.Exit(2)
	}
	return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
}

// Generate makes a certificate (e.g., key and secret)
// If they are ephemeral, I think we can do this
func (c *Certificate) Generate() {

	// Note - we could use rsa here too.
	private, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Flux Framework"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, host := range c.Hosts {

		// Does the host have an ip address?
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

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &private.PublicKey, private)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	c.Public = out.String()
	fmt.Println(c.Public)
	out.Reset()
	pem.Encode(out, pemBlockForKey(private))
	c.Private = out.String()
	fmt.Println(c.Private)
}
