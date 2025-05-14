package main

import (
	"fmt"
	"log"
	"net/http"
	"gokube/cmd/gokube/config"
)


//healthHandler is a function which returns the health status of the service.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "healthy")
}



func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}



func main() {

	config := config.LoadAppConfig()
	fmt.Printf("Loaded config: %+v\n", config)
	// Register handlers
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", helloHandler)

	// Start server
	log.Printf("Server starting on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

