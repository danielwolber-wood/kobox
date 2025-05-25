package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	refreshOptions, err := AuthFlow()
	if err != nil {
		log.Fatal("error with auth flow: %v\n", err)
	}

	s, err := newServer(*refreshOptions)
	if err != nil {
		log.Fatalf("error creating server: %v\n", err.Error())
	}

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go s.worker(i)
	}

	r := http.NewServeMux()
	r.HandleFunc("/health", handleHealthCheck)
	r.HandleFunc("/v2/api/upload/url", s.handlerUploadURL)
	r.HandleFunc("/v2/api/upload/html", s.handlerUploadFullPage)
	fmt.Println("Serving on :8080")
	err = http.ListenAndServeTLS(":8080", "server.crt", "server.key", r) // I tried this line and everything returned 404
	//err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
