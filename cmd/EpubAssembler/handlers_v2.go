package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/danielwolber-wood/kobox/internal/response"
	"io"
	"net/http"
)

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("server is alive"))
}

func (s *Server) handlerUploadURL(w http.ResponseWriter, r *http.Request) {
	var obj URLRequestObject
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("error reading body:%v\n", err))
		return
	}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, fmt.Sprintf("error parsing request:%v\n", err))
		return
	}
	defer r.Body.Close()
	url := obj.Url
	job := Job{taskType: TaskFetch, url: url, generateOptions: GenerateOptions{Title: obj.Title}}
	s.jobQueue <- job
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handlerUploadFullPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set CORS headers for actual request
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var obj HTMLRequestObject
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("error reading body:%v\n", err))
		return
	}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, fmt.Sprintf("error parsing request:%v\n", err))
		return
	}
	defer r.Body.Close()
	html := obj.Html
	job := Job{taskType: TaskExtract, htmlReader: bytes.NewReader([]byte(html)), fullText: html, generateOptions: GenerateOptions{Title: obj.Title}}
	s.jobQueue <- job
	w.WriteHeader(http.StatusOK)
}
