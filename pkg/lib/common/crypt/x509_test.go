package crypt

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestVerifyX509CertChainGoodChain(t *testing.T) {
	rootCAPkixName := pkix.Name{
		CommonName:    "Acme Corp Signing Root CA",
		Organization:  []string{"Acme"},
		Country:       []string{"US"},
		Province:      []string{"CA"},
		Locality:      []string{"Santa Clara"},
		StreetAddress: []string{"123 Anony Mouse Blvd."},
		PostalCode:    []string{"12345"},
	}

	intermediate1PkixName := pkix.Name{
		CommonName: "Acme TPM Intermediate CA",
	}

	ekCertPkixName := pkix.Name{
		CommonName: "Acme TPM EK Cert",
	}

	// Generate a self signed root CA
	caPrivateKey, caPubkey, _ := GenerateKeyPair("rsa", 4096)

	rootCaTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(2020),
		Subject:               rootCAPkixName,
		NotBefore:             time.Now().AddDate(-1, 0, 0),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Create intermediate Certs for signing the leaf
	intermediateCert1Template := x509.Certificate{
		SerialNumber:          big.NewInt(2021),
		Subject:               intermediate1PkixName,
		NotBefore:             time.Now().AddDate(-1, 0, 0),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Create the chain starting with Root
	rootCertBytes, err := x509.CreateCertificate(rand.Reader, &rootCaTemplate, &rootCaTemplate,
		caPubkey, caPrivateKey)
	rootCertx509, err := x509.ParseCertificate(rootCertBytes)

	// INTER 1
	intermediate1CertBytes, err := x509.CreateCertificate(rand.Reader, &intermediateCert1Template, rootCertx509,
		caPubkey, caPrivateKey)
	intermediate1Certx509, err := x509.ParseCertificate(intermediate1CertBytes)

	// LEAF
	ekCertTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2023),
		Subject:      ekCertPkixName,
		NotBefore:    time.Now().AddDate(-2, 0, 0),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		KeyUsage:     x509.KeyUsageEncipherOnly,
	}

	// create the EK leaf certificate
	ekCertificateBytes, err := x509.CreateCertificate(rand.Reader, &ekCertTemplate, intermediate1Certx509,
		caPubkey, caPrivateKey)
	t.Log(err)

	ekCertx509, err := x509.ParseCertificate(ekCertificateBytes)
	t.Log(err)

	// combine all certs
	var allCerts []*x509.Certificate
	allCerts = append(allCerts, rootCertx509, intermediate1Certx509, ekCertx509)

	assert.NoError(t, VerifyX509CertChain(true, allCerts, nil))
}

func TestVerifyX509CertChainExpired(t *testing.T) {
	rootCAPkixName := pkix.Name{
		CommonName:    "Acme Corp Signing Root CA",
		Organization:  []string{"Acme"},
		Country:       []string{"US"},
		Province:      []string{"CA"},
		Locality:      []string{"Santa Clara"},
		StreetAddress: []string{"123 Anony Mouse Blvd."},
		PostalCode:    []string{"12345"},
	}

	intermediate1PkixName := pkix.Name{
		CommonName: "Acme TPM Model CA",
	}

	ekCertPkixName := pkix.Name{
		CommonName: "Acme TPM EK Cert",
	}

	// Generate a self signed root CA
	caPrivateKey, caPubkey, err := GenerateKeyPair("rsa", 4096)

	rootCaTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(2020),
		Subject:               rootCAPkixName,
		NotBefore:             time.Now().AddDate(-1, 0, 0),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create intermediate Certs for signing the leaf
	intermediateCert1Template := x509.Certificate{
		SerialNumber:          big.NewInt(2021),
		Subject:               intermediate1PkixName,
		NotBefore:             time.Now().AddDate(-1, 0, 0),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create the chain starting with Root
	rootCertBytes, err := x509.CreateCertificate(rand.Reader, &rootCaTemplate, &rootCaTemplate, caPubkey, caPrivateKey)
	rootCertx509, err := x509.ParseCertificate(rootCertBytes)

	// INTER 1
	intermediate1CertBytes, err := x509.CreateCertificate(rand.Reader, &intermediateCert1Template, rootCertx509,
		rootCertx509.PublicKey, caPrivateKey)
	intermediate1Certx509, err := x509.ParseCertificate(intermediate1CertBytes)

	leafPrivKey, leafPubKey, err := GenerateKeyPair("rsa", 4096)

	// LEAF
	ekCertTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(2023),
		Subject:               ekCertPkixName,
		NotBefore:             time.Now().AddDate(-2, 0, 0),
		NotAfter:              time.Now().AddDate(-1, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
		KeyUsage:              x509.KeyUsageEncipherOnly,
		OCSPServer:            []string{"http://ocsp.example.com"},
		IssuingCertificateURL: []string{"http://crt.example.com/ca1.crt"},
		DNSNames:              []string{"test.example.com"},
		EmailAddresses:        []string{"somebody@thatiusedtoknow.org"},
		ExtraExtensions: []pkix.Extension{
			{
				Id:    []int{1, 2, 3, 4},
				Value: []byte("extra extension"),
			},
			// This extension should override the SubjectKeyId, above.
			{
				Id:       []int{2, 5, 29, 14},
				Critical: false,
				Value:    []byte{0x04, 0x04, 4, 3, 2, 1},
			},
		},
	}

	// create the EK leaf certificate
	ekCertificateBytes, err := x509.CreateCertificate(rand.Reader, &ekCertTemplate, &intermediateCert1Template,
		leafPubKey, leafPrivKey)

	ekCertx509, err := x509.ParseCertificate(ekCertificateBytes)
	t.Log(err)
	// combine all certs
	var allCerts []*x509.Certificate
	allCerts = append(allCerts, rootCertx509, intermediate1Certx509, ekCertx509)

	assert.Error(t, VerifyX509CertChain(true, allCerts, nil))
}
