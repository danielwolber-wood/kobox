package main

import (
	"log"
	"os"
)

type Server struct {
	dropboxApiKey string
	jobQueue      chan Job
}

func newServer() (*Server, error) {
	apiKey := os.Getenv("DROPBOX_API_KEY")
	if apiKey == "" {
		log.Fatal("DROPBOX_API_KEY not set")
	}
	jobQueue := make(chan Job, 256)
	return &Server{
		dropboxApiKey: apiKey,
		jobQueue:      jobQueue,
	}, nil
}
