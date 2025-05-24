package main

import (
	"log"
	"os"
)

type Server struct {
	readabilityParser *ReadabilityParser
	dropboxApiKey     string
}

func newServer() (*Server, error) {
	apiKey := os.Getenv("DROPBOX_API_KEY")
	if apiKey == "" {
		log.Fatal("DROPBOX_API_KEY not set")
	}
	parser, err := NewReadabilityParser()
	if err != nil {
		return nil, err
	}
	return &Server{
		readabilityParser: parser,
		dropboxApiKey:     apiKey,
	}, nil
}
