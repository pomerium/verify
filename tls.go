package verify

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type tlsVerifier struct {
	mu    sync.Mutex
	valid map[string]struct{}
}

func newTLSVerifier() *tlsVerifier {
	return &tlsVerifier{valid: make(map[string]struct{})}
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
		v.valid[serverName] = struct{}{}
	} else {
		delete(v.valid, serverName)
	}
	v.mu.Unlock()
	return nil
}

func (v *tlsVerifier) IsValid(r *http.Request) bool {
	if r == nil || r.TLS == nil {
		return false
	}

	v.mu.Lock()
	_, ok := v.valid[r.TLS.ServerName]
	v.mu.Unlock()

	return ok
}

func tlsHost(targetAddr string) string {
	if strings.LastIndex(targetAddr, ":") > strings.LastIndex(targetAddr, "]") {
		targetAddr = targetAddr[:strings.LastIndex(targetAddr, ":")]
	}
	return targetAddr
}
