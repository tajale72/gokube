package router

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

const message = "Hello World!!"

type Handlers struct {
	logger *log.Logger
}

func (h *Handlers) SetupRouttes(mux *http.ServeMux) {
	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Handle the root route
	mux.HandleFunc("/", h.Logger(h.HomeHandler))
	mux.HandleFunc("/health", h.Logger(h.healthHandler))

}

func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the HTML file for the root path
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "static/index.html")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

// healthHandler returns the health status of the service.
func (h *Handlers) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s %s", runtime.Version(), "v1.0.0")
}

func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("request processed in %s\n", time.Since(startTime))
		next(w, r)

	}
}
