package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"sync"
	"time"
)

var certMap = make(map[string]*tls.Certificate)
var mapLock = &sync.RWMutex{}

func GetOrCreateCertificate(host string) (*tls.Certificate, error) {
	mapLock.RLock()
	cert, exists := certMap[host]
	mapLock.RUnlock()

	if exists {
		return cert, nil
	}

	// Certificate does not exist, generate a new one
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Create a template for the certificate
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Self-Signed"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 1 year
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create a self-signed certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, err
	}

	// Create a TLS certificate using the private key and the certificate bytes
	tlsCert := tls.Certificate{
		Certificate: [][]byte{certBytes},
		PrivateKey:  privKey,
	}

	// Store the new certificate in the map
	mapLock.Lock()
	certMap[host] = &tlsCert
	mapLock.Unlock()

	return &tlsCert, nil
}
