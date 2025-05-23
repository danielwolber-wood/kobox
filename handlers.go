package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("server is alive"))
}

func (s *Server) handleReadabilityURL(w http.ResponseWriter, r *http.Request) {
	// parse json from response
	// get sought URL
	// get URL data with http.Get()
	// send to readability module
	// return RedabilityResult
	var req ReadabilityURLRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Sprintf("cannot parse json: %v\n", err))
		return
	}

	resp, err := http.Get(req.Url)
	if err != nil {
		WriteError(w, http.StatusBadGateway, fmt.Sprintf("cannot get page: %v\n", err))
		return
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("cannot parse get request: %v\n", err))
		return
	}

	result, err := s.readabilityParser.ParseHTML(string(body))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("readability error: %v\n", err))
		return
	}

	WriteJsonResponse(w, http.StatusOK, result)

}
