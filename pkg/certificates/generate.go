package certificates

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"jaeger-tenant/pkg/structures"
	"math/big"
	mathRand "math/rand"
	"time"
)

// Certificate template
type Certificate struct {
	CommonName         string
	Locality           []string
	Organization       []string
	OrganizationalUnit []string
	Country            []string
	DNSNames           []string
}

func generateCAKeyPair(crtTemplate *Certificate) (caCert, caKey *bytes.Buffer) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	// Generate a pem block with the private key
	/*keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})*/
	keyPem := bytes.NewBufferString("")
	_ = pem.Encode(keyPem, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	tml := x509.Certificate{
		// you can add any attr that you need
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(5, 0, 0),
		// you have to generate a different serial number each execution
		SerialNumber: big.NewInt(123123),
		Subject: pkix.Name{
			CommonName:         crtTemplate.CommonName,
			Organization:       crtTemplate.Organization,
			OrganizationalUnit: crtTemplate.OrganizationalUnit,
			Country:            crtTemplate.Country,
		},
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	cert, _ := x509.CreateCertificate(rand.Reader, &tml, &tml, &key.PublicKey, key)

	// Generate a pem block with the certificate
	certPem := bytes.NewBufferString("")
	_ = pem.Encode(certPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	/*
		tlsCert, err := tls.X509KeyPair(certPem.Bytes(), keyPem.Bytes())
		if err != nil {
			log.Fatal("Cannot be loaded the certificate.", err.Error())
		}
	*/

	return certPem, keyPem
}

func generateCSRKeyPair(crtTemplate *Certificate) (clientCSR, clientKey *bytes.Buffer) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	// Generate a pem block with the private key
	/*keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})*/
	clientKey = bytes.NewBufferString("")
	_ = pem.Encode(clientKey, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	subj := pkix.Name{
		CommonName:         crtTemplate.CommonName,
		Country:            crtTemplate.Country,
		Locality:           crtTemplate.Locality,
		Organization:       crtTemplate.Organization,
		OrganizationalUnit: crtTemplate.OrganizationalUnit,
	}
	rawSubj := subj.ToRDNSequence()
	/*
		var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}
		emailAddress := "arpit.agarwal02@sap.com"
		rawSubj = append(rawSubj, []pkix.AttributeTypeAndValue{
			{Type: oidEmailAddress, Value: emailAddress},
		})
	*/

	asn1Subj, _ := asn1.Marshal(rawSubj)
	template := x509.CertificateRequest{
		RawSubject: asn1Subj,
		//EmailAddresses:     []string{emailAddress},
		SignatureAlgorithm: x509.SHA256WithRSA,
		DNSNames:           crtTemplate.DNSNames,
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, key)
	clientCSR = bytes.NewBufferString("")
	_ = pem.Encode(clientCSR, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	return
}

func generateClientCRT(caCRT, caKey, clientCSR *bytes.Buffer) (clientCRT *bytes.Buffer, err error) {

	block, _ := pem.Decode(caCRT.Bytes())
	parsedCaCrt, _ := x509.ParseCertificate(block.Bytes)

	block, _ = pem.Decode(caKey.Bytes())
	parsedCaKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	block, _ = pem.Decode(clientCSR.Bytes())
	parsedClientCsr, _ := x509.ParseCertificateRequest(block.Bytes)

	tml := x509.Certificate{
		// you can add any attr that you need
		Raw:                     parsedClientCsr.Raw,
		RawTBSCertificate:       parsedClientCsr.RawTBSCertificateRequest,
		RawSubjectPublicKeyInfo: parsedClientCsr.RawSubjectPublicKeyInfo,
		RawSubject:              parsedClientCsr.RawSubject,
		Signature:               parsedClientCsr.Signature,
		SignatureAlgorithm:      parsedClientCsr.SignatureAlgorithm,
		PublicKeyAlgorithm:      parsedClientCsr.PublicKeyAlgorithm,
		PublicKey:               parsedClientCsr.PublicKey,
		Version:                 parsedClientCsr.Version,
		Subject:                 parsedClientCsr.Subject,
		Extensions:              parsedClientCsr.Extensions,
		ExtraExtensions:         parsedClientCsr.ExtraExtensions,
		DNSNames:                parsedClientCsr.DNSNames,
		EmailAddresses:          parsedClientCsr.EmailAddresses,
		IPAddresses:             parsedClientCsr.IPAddresses,
		URIs:                    parsedClientCsr.URIs,
		NotBefore:               time.Now(),
		NotAfter:                time.Now().AddDate(5, 0, 0),
		SerialNumber:            big.NewInt(1),
	}

	clientCrtBytes, err := x509.CreateCertificate(rand.Reader, &tml, parsedCaCrt, parsedClientCsr.PublicKey, parsedCaKey)
	if err != nil {
		return
	}
	clientCRT = bytes.NewBufferString("")
	_ = pem.Encode(clientCRT, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCrtBytes,
	})
	return
}

// GenerateCredentials fills tenant with credentials
func GenerateCredentials(tenant *structures.Tenant) {

	crtTemplate := &Certificate{
		CommonName:         tenant.Customer,
		Locality:           []string{"Bangalore"},
		Organization:       []string{"SAP CP Core"},
		OrganizationalUnit: []string{"Jaeger Service Provider"},
		Country:            []string{"IN"},
		DNSNames:           []string{tenant.JaegerCollectorURL},
	}

	caCrt, caKey := generateCAKeyPair(crtTemplate)
	clientCsr, clientKey := generateCSRKeyPair(crtTemplate)
	clientCrt, _ := generateClientCRT(caCrt, caKey, clientCsr)
	tenant.CACrt = caCrt.String()
	tenant.CAKey = caKey.String()
	tenant.ClientCrt = clientCrt.String()
	tenant.ClientKey = clientKey.String()
	tenant.Username, tenant.Password = generateBasicAuthCred(10)
}

func generateBasicAuthCred(strLen int) (username, password string) {
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789}{$#!@"
	result := make([]byte, strLen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	username = "jaeger"
	password = string(result)
	return
}
