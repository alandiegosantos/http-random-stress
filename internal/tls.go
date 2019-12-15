package internal

import (
	"crypto/tls"
	"crypto/x509"
)

func GenerateTlsConfig() *tls.Config {

	// Create a pool with the server certificate since it is not signed
	// by a known CA
	caCertPool := x509.NewCertPool()

	// Create TLS configuration with the certificate of the server
	return &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            caCertPool,
	}

}
