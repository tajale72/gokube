package router

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestHandlers() *Handlers {
	return &Handlers{
		logger: log.New(io.Discard, "", 0),
	}
}

func TestNewHandlers(t *testing.T) {
	logger := log.New(io.Discard, "", 0)

	h := NewHandlers(logger)

	if h == nil {
		t.Fatal("expected handlers but got nil")
	}

	if h.logger != logger {
		t.Fatal("expected logger to be set")
	}
}

func TestHealthHandler(t *testing.T) {
	h := newTestHandlers()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	h.HealthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	body := rr.Body.String()

	if !strings.Contains(body, "go") {
		t.Fatalf("expected Go runtime version in body, got %q", body)
	}

	if !strings.Contains(body, "v1.0.0") {
		t.Fatalf("expected app version in body, got %q", body)
	}
}

func TestEnforceAuth_MissingAPIKey(t *testing.T) {
	t.Setenv("X_API_KEY", "secret")

	h := newTestHandlers()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	rr := httptest.NewRecorder()

	h.EnforceAuth(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d but got %d", http.StatusUnauthorized, rr.Code)
	}

	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}
}

func TestEnforceAuth_InvalidAPIKey(t *testing.T) {
	t.Setenv("X_API_KEY", "secret")

	h := newTestHandlers()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	req.Header.Set("X-API-KEY", "wrong-key")

	rr := httptest.NewRecorder()

	h.EnforceAuth(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d but got %d", http.StatusUnauthorized, rr.Code)
	}

	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}
}

func TestEnforceAuth_ValidAPIKey(t *testing.T) {
	t.Setenv("X_API_KEY", "secret")

	h := newTestHandlers()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	req.Header.Set("X-API-KEY", "secret")

	rr := httptest.NewRecorder()

	h.EnforceAuth(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}
}

func TestLoggerMiddleware(t *testing.T) {
	h := newTestHandlers()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusAccepted)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	h.Logger(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected status %d but got %d", http.StatusAccepted, rr.Code)
	}

	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}
}

func TestSetupRoutes(t *testing.T) {
	h := newTestHandlers()

	mux := http.NewServeMux()
	h.SetupRoutes(mux)

	t.Setenv("X_API_KEY", "secret")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected health route status %d but got %d", http.StatusOK, rr.Code)
	}
}

func TestHomeHandler_IndexFile(t *testing.T) {
	h := newTestHandlers()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	h.HomeHandler(rr, req)

	// Could be 200 if file exists
	// Could be 404 during CI/tests if static file missing
	if rr.Code != http.StatusOK && rr.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d or %d but got %d",
			http.StatusOK,
			http.StatusNotFound,
			rr.Code,
		)
	}
}

func TestHomeHandler_StaticAsset(t *testing.T) {
	h := newTestHandlers()

	req := httptest.NewRequest(
		http.MethodGet,
		"/style.css",
		nil,
	)

	rr := httptest.NewRecorder()

	h.HomeHandler(rr, req)

	// Could be 200 if asset exists
	// Could be 404 if asset missing
	if rr.Code != http.StatusOK && rr.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d or %d but got %d",
			http.StatusOK,
			http.StatusNotFound,
			rr.Code,
		)
	}
}
