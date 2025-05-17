package server

import (
	"crypto/tls"
	"net/http"
	"time"
)

func NewServer(mux *http.ServeMux, serverAddr string, serverPort string) *http.Server {

	//Initialize the TLS configuration
	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}

	srv := &http.Server{
		Addr:         serverAddr + ":" + serverPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
		TLSConfig:    tlsConfig,
	}
	return srv
}
