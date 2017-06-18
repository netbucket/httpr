// Copyright © 2017 Igor Bondarenko <ibondare@protonmail.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"time"
)

const (
	rsaKeyLength     = 2048
	serialNumberBits = 128
)

// StartHTTPSListener starts an HTTPS server at the address specified by the service parameter
// If either or both certFile and keyFile are blank, a self-singned cert is generated
func StartHTTPSListener(service, certFile, keyFile string) error {
	s := http.Server{}

	// If certFile and/or keyFile are blank, generate a self-signed TLS cert
	if len(certFile) == 0 || len(keyFile) == 0 {
		selfSignedCert, err := generateSelfSignedCert()

		if err != nil {
			return err
		}

		s.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{selfSignedCert},
		}
	}

	s.Addr = service

	return s.ListenAndServeTLS(certFile, keyFile)
}

// Generate a self-signed TLS cert for the HTTPS endpoint
func generateSelfSignedCert() (tls.Certificate, error) {
	rootKey, err := rsa.GenerateKey(rand.Reader, rsaKeyLength)
	if err != nil {
		return tls.Certificate{}, err
	}

	t, err := createX509Template()
	if err != nil {
		return tls.Certificate{}, err
	}

	t.IsCA = true
	t.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	t.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	t.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}

	rootCertPEM, err := createCertFromTemplate(t, rootKey)

	if err != nil {
		return tls.Certificate{}, err
	}

	// Print the self-signed cert
	fmt.Printf("%s\n", rootCertPEM)

	// PEM encode the private key
	rootKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rootKey),
	})

	// Create a TLS cert using the private key and certificate
	return tls.X509KeyPair(rootCertPEM, rootKeyPEM)
}

// Create a certificate template
func createX509Template() (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), serialNumberBits)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return nil, err
	}

	t := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"HTTP Rake - httpr"}},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365), // Make it valid for a year
		BasicConstraintsValid: true,
	}

	return &t, nil
}

// Create a self-signed certificate, PEM-encoded in an in-memory byte array, using a supplied template
func createCertFromTemplate(template *x509.Certificate, key *rsa.PrivateKey) (certPEM []byte, err error) {
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return
	}

	_, err = x509.ParseCertificate(certDER)
	if err != nil {
		return
	}

	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM = pem.EncodeToMemory(&b)

	return
}
