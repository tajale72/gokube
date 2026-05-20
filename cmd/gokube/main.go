package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"gokube/config"
	"gokube/router"
	"gokube/server"
)

var (
	// myCert = os.Getenv("MY_CERT")
	// myKey  = os.Getenv("MY_KEY")
	// addr = os.Getenv("ADDR")
	// port = os.Getenv("PORT")

	myCert = "server.crt"
	myKey  = "server.key"
	addr   = "localhost"
	port   = "8080"
)

const (
	serviceName = "dragon"
)

func main() {
	config := config.LoadAppConfig()
	fmt.Printf("Loaded config: %+v\n", config)

	logger := log.New(os.Stdout, serviceName+" ", log.LstdFlags|log.Lshortfile)

	h := router.NewHandlers(logger)

	mux := http.NewServeMux()
	h.SetupRoutes(mux)

	// Wrap mux with CORS middleware
	handlerWithRateLimit := rateLimiterMiddleware(mux)
	handlerWithCors := corsMiddleware(handlerWithRateLimit)

	srv := server.NewServer(handlerWithCors, addr, port)

	logger.Printf("Starting server at https://%s:%s\n", addr, port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v ", err)
	}
}

func rateLimiterMiddleware(next http.Handler) http.Handler {
	type visitor struct {
		count     int
		resetTime time.Time
	}

	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		now := time.Now()

		mu.Lock()

		v, exists := visitors[ip]
		if !exists || now.After(v.resetTime) {
			visitors[ip] = &visitor{
				count:     1,
				resetTime: now.Add(time.Minute),
			}
			mu.Unlock()

			next.ServeHTTP(w, r)
			return
		}

		if v.count >= 60 {
			mu.Unlock()
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		v.count++
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
