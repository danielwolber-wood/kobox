package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	s, err := newServer()
	if err != nil {
		log.Fatalf("error creating server: %v\n", err.Error())
	}
	r := http.NewServeMux()
	r.HandleFunc("/health", handleHealthCheck)
	r.HandleFunc("/v1/api/readability", s.handleReadabilityURL)
	r.HandleFunc("/v1/api/assembler", s.handleAssembler)
	//r.HandleFunc("/v2/api/fetch", s.handleFetch)
	//r.HandleFunc("/v2/api/extract", s.handleExtract)
	//r.HandleFunc("/v2/api/generate", s.handleGenerate)
	//r.HandleFunc("/v2/api/upload", s.handleUpload)
	fmt.Println("Serving on :8080")
	http.ListenAndServe(":8080", r)

}
