package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	serviceName = "kafkaService"
)

func main() {

	config := config.LoadAppConfig()
	fmt.Printf("Loaded config: %+v\n", config)

	//Initializing the loggervb
	logger := log.New(os.Stdout, serviceName+" ", log.LstdFlags|log.Lshortfile)

	h := router.NewHandlers(logger)

	//Initialize the server
	mux := http.NewServeMux()

	//
	h.SetupRouttes(mux)

	//Initializing the server
	srv := server.NewServer(mux, addr, port)

	// log starting the server
	logger.Printf("Starting server at https://%s:%s\n", addr, port)

	//Listen and serve
	err := srv.ListenAndServeTLS(myCert, myKey)
	if err != nil {
		log.Fatalf("Server failed to start: %v ", err)
	}
}
