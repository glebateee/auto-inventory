package options

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

const (
	srvCert = "./internal/grpc/server/certs/server.crt"
	srvKey  = "./internal/grpc/server/certs/server.key"
	caCert  = "./internal/grpc/server/certs/ca.crt"
)

func NewTLS() (credentials.TransportCredentials, error) {
	const op = "options.NewTLS"
	cert, err := tls.LoadX509KeyPair(srvCert, srvKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(caCert)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("%s: %w", op, errors.New("failed to append ca certificate"))
	}
	return credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS12,
	}), nil
}
