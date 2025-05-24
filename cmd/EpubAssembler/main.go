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
	//r.HandleFunc("/v2/api/url-to-rr", s.handleUrlToRR)
	//r.HandleFunc("/v2/api/url-to-epub", s.handleUrlToEpub)
	//r.HandleFunc("/v2/api/url-to-dropbox", s.handleUrlToDropbox)
	//r.HandleFunc("/v2/api/page-to-rr", s.handlePageToRR) // same as extract
	//r.HandleFunc("/v2/api/page-to-epub", s.handlePageToEpub)
	//r.HandleFunc("/v2/api/page-to-dropbox", s.handlePageToDropbox)
	//r.HandleFunc("/v2/api/rr-to-epub", s.handleRRToEpub) // same as generate
	//r.HandleFunc("/v2/api/rr-to-dropbox", s.handleRRToDropbox)
	//r.HandleFunc("/v2/api/epub-to-dropbox", s.handleEpubToDropbox) // same as upload
	fmt.Println("Serving on :8080")
	http.ListenAndServe(":8080", r)

}
