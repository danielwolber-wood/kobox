package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	refreshOptions, err := AuthFlowPKCE()
	if err != nil {
		log.Fatalf("error with auth flow: %v\n", err)
	}

	s, err := newServer(*refreshOptions)
	if err != nil {
		log.Fatalf("error creating server: %v\n", err.Error())
	}

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go s.worker(i)
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Print("cannot load environment file")
	}
	port := os.Getenv("KOBOX_PORT")
	r := http.NewServeMux()
	r.HandleFunc("/health", handleHealthCheck)
	r.HandleFunc("/v2/api/upload/url", s.handlerUploadURL)
	r.HandleFunc("/v2/api/upload/html", s.handlerUploadFullPage)
	fmt.Println("Serving on %v\n", port)
	err = http.ListenAndServeTLS(port, "server.crt", "server.key", r)
	if err != nil {
		panic(err)
	}
}
