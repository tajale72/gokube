package router

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"gokube/config"
)

const message = "Hello World!!"

type Handlers struct {
	logger *log.Logger
}

func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	mux.Handle("/", h.Logger(http.HandlerFunc(h.HomeHandler)))
	mux.Handle("/health", h.Logger(http.HandlerFunc(h.HealthHandler)))

	// Protected weather route remains perfect
	protectedWeather := h.EnforceAuth(http.HandlerFunc(h.WeatherHandler))
	mux.Handle("/weather", h.Logger(protectedWeather))
}

func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {

	h.logger.Println("calling handler")
	// If the user requests a specific static file asset (like script.js or style.css)
	if r.URL.Path != "/" {
		http.ServeFile(w, r, "./static"+r.URL.Path)
		return
	}

	// Default to serving index.html for the root route
	http.ServeFile(w, r, "./static/index.html")
}

// healthHandler returns the health status of the service.
func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s %s", runtime.Version(), "v1.0.0")
}

func (h *Handlers) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("request processed in %s\n", time.Since(startTime))
		next.ServeHTTP(w, r)

	})
}

// enforceAuth checks for a valid X-API-Key header before allowing access
func (h *Handlers) EnforceAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")

		// Fail fast if the token is missing or incorrect
		if apiKey == "" || apiKey != config.GetString("X_API_KEY", "deafult") {
			h.logger.Printf("Invalid or missing X-API-KEY %s\n", apiKey)
			http.Error(w, "Unauthorized: Invalid or missing X-API-KEY", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the actual endpoint
		next.ServeHTTP(w, r)
	})
}
