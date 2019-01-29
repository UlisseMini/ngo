package tlsconfig

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"
)

// Get creates a self signed certificate and
// returns a valid *tls.Config for listening.
func Get(host string) (*tls.Config, error) {
	if len(host) == 0 {
		return nil, fmt.Errorf("invalid host: %q", host)
	}

	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %s", err)
	}

	// Manage the time
	notBefore := time.Now()
	notAfter := notBefore.Add(30 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	pub := &priv.PublicKey

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template,
		pub, priv)

	if err != nil {
		return nil, err
	}

	// Encode the raw bytes into PEM
	keyBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	keyPEMBlock := &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}
	certPEMBlock := &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}

	// Convert pemCert into a tls.Certificate
	tlsCert, err := tls.X509KeyPair(pem.EncodeToMemory(certPEMBlock),
		pem.EncodeToMemory(keyPEMBlock))
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}, nil
}
