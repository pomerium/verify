package verify

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type tlsVerifier struct {
	mu     sync.Mutex
	errors map[string]error
}

func newTLSVerifier() *tlsVerifier {
	return &tlsVerifier{errors: make(map[string]error)}
}

func (v *tlsVerifier) DialTLSContext(ctx context.Context, network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	log.Info().Str("addr", addr).Msg("dialing")
	conn, err := tls.DialWithDialer(dialer, network, addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         tlsHost(addr),
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return v.VerifyPeerCertificate(tlsHost(addr), rawCerts)
		},
	})
	if err != nil {
		log.Error().Err(err).Str("addr", addr).Msg("error dialing TLS")
		return nil, err
	}
	return conn, nil
}

func (v *tlsVerifier) VerifyPeerCertificate(serverName string, rawCerts [][]byte) error {
	certs := make([]*x509.Certificate, len(rawCerts))
	for i, rawCert := range rawCerts {
		var err error
		certs[i], err = x509.ParseCertificate(rawCert)
		if err != nil {
			log.Error().Err(err).Msg("error parsing TLS certificate")
			return err
		}
	}

	opts := x509.VerifyOptions{
		DNSName:       serverName,
		Intermediates: x509.NewCertPool(),
	}
	for _, cert := range certs[1:] {
		opts.Intermediates.AddCert(cert)
	}
	_, err := certs[0].Verify(opts)
	v.mu.Lock()
	if err == nil {
		delete(v.errors, serverName)
	} else {
		v.errors[serverName] = err
		log.Error().
			Err(err).
			Str("server-name", serverName).
			Msg("invalid TLS certificate")
	}
	v.mu.Unlock()
	return nil
}

func (v *tlsVerifier) GetTLSError(serverName string) error {
	v.mu.Lock()
	err := v.errors[serverName]
	v.mu.Unlock()
	return err
}

func tlsHost(targetAddr string) string {
	if strings.LastIndex(targetAddr, ":") > strings.LastIndex(targetAddr, "]") {
		targetAddr = targetAddr[:strings.LastIndex(targetAddr, ":")]
	}
	return targetAddr
}
