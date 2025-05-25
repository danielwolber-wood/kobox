package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	/*
		refreshOptions, err := AuthFlowPKCE()
		if err != nil {
			log.Fatalf("error with auth flow: %v\n", err)
		}
	*/

	ensureCertificates()
	s, err := newServer()
	if err != nil {
		log.Fatalf("error creating server: %v\n", err.Error())
	}

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go s.worker(i)
	}

	err = godotenv.Load(envFile)
	if err != nil {
		log.Print("cannot load environment file")
	}
	port := os.Getenv("KOBOX_PORT")
	if port == "" {
		port = ":12332"
	}
	r := http.NewServeMux()
	r.HandleFunc("/health", handleHealthCheck)
	r.HandleFunc("/config", s.handlerConfig)
	r.HandleFunc("/v2/api/upload/url", s.handlerUploadURL)
	r.HandleFunc("/v2/api/upload/html", s.handlerUploadFullPage)
	fmt.Printf("Serving on %v\n", port)
	err = http.ListenAndServeTLS(port, "/app/certs/server.crt", "/app/certs/server.key", r)
	if err != nil {
		panic(err)
	}
}
