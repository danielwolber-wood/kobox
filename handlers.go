package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("server is alive"))
}

func (s *Server) handleReadabilityURL(w http.ResponseWriter, r *http.Request) {
	var req ReadabilityURLRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Sprintf("cannot parse json: %v\n", err))
		return
	}
	result, err := s.readabilityParser.ParseURL(req.Url)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("cannot get readability: %v\n", err))
		return
	}
	WriteJsonResponse(w, http.StatusOK, result)
}

func (s *Server) handleAssembler(w http.ResponseWriter, r *http.Request) {
	// accepts a readabilityResult JSON and returns epub
	var readabilityResult ReadabilityResult
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&readabilityResult)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Sprintf("cannot parse json: %v\n", err))
		return
	}
	// accepts a ReadabilityResult and returns a .EPUB
	html := GenerateHTML(readabilityResult.Title, readabilityResult.Content)
	epubData, err := ConvertStringWithPandoc(html, "html", "epub")
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("cannot convert file: %v\n", err))
		return
	}
	WriteEpubResponse(w, http.StatusOK, epubData, readabilityResult.Title)
}
