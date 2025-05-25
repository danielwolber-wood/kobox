package main

import (
	"encoding/json"
	"fmt"
	"github.com/danielwolber-wood/kobox/internal/response"
	"io"
	"net/http"
)

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
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
	job := Job{TaskFetch, url, "", nil, GenerateOptions{Title: obj.Title}, UploadOptions{}}
	s.jobQueue <- job
}

func handlerUploadFullPage(w http.ResponseWriter, r *http.Request) {
}

func handlerUploadReadabilityObject(w http.ResponseWriter, r *http.Request) {}

func handlerUploadEpub(w http.Request, r *http.Request) {}
