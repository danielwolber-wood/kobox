package main

import (
	"log"
	"os"
)

type Server struct {
	dropboxApiKey string
	FetchQueue    chan URL
	ExtractQueue  chan HTML
	GenerateQueue chan ReadabilityObject
	UploadQueue   chan UploadObject
}

func newServer() (*Server, error) {
	apiKey := os.Getenv("DROPBOX_API_KEY")
	if apiKey == "" {
		log.Fatal("DROPBOX_API_KEY not set")
	}
	fetchQueue := make(chan URL, 16)
	extractQueue := make(chan HTML, 16)
	generateQueue := make(chan ReadabilityObject, 16)
	uploadQueue := make(chan UploadObject, 32)
	return &Server{
		dropboxApiKey: apiKey,
		FetchQueue:    fetchQueue,
		ExtractQueue:  extractQueue,
		GenerateQueue: generateQueue,
		UploadQueue:   uploadQueue,
	}, nil
}
