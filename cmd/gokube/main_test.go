package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsMiddleware_OptionsRequest(t *testing.T) {
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	handler := corsMiddleware(next)

	req := httptest.NewRequest(http.MethodOptions, "/weather", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d but got %d", http.StatusNoContent, rr.Code)
	}

	if nextCalled {
		t.Fatal("expected next handler not to be called for OPTIONS request")
	}

	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("expected CORS origin header")
	}
}

func TestCorsMiddleware_NormalRequest(t *testing.T) {
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := corsMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}

	if rr.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Fatal("expected CORS methods header")
	}
}
