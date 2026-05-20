package server

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	mux := http.NewServeMux()

	srv := NewServer(mux, "localhost", "8080")

	if srv.Addr != "localhost:8080" {
		t.Fatalf(
			"expected addr %q but got %q",
			"localhost:8080",
			srv.Addr,
		)
	}

	if srv.Handler != mux {
		t.Fatal("expected handler to match mux")
	}

	if srv.ReadTimeout != 10*time.Second {
		t.Fatalf(
			"expected read timeout %v but got %v",
			10*time.Second,
			srv.ReadTimeout,
		)
	}

	if srv.WriteTimeout != 10*time.Second {
		t.Fatalf(
			"expected write timeout %v but got %v",
			10*time.Second,
			srv.WriteTimeout,
		)
	}

	if srv.IdleTimeout != 120*time.Second {
		t.Fatalf(
			"expected idle timeout %v but got %v",
			120*time.Second,
			srv.IdleTimeout,
		)
	}

	if srv.TLSConfig == nil {
		t.Fatal("expected TLS config but got nil")
	}

	if srv.TLSConfig.MinVersion != tls.VersionTLS13 {
		t.Fatalf(
			"expected TLS version %v but got %v",
			tls.VersionTLS13,
			srv.TLSConfig.MinVersion,
		)
	}

	if len(srv.TLSConfig.CipherSuites) == 0 {
		t.Fatal("expected cipher suites but got none")
	}
}
