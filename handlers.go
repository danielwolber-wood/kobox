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

func (s *Server) handleFeedPost(w http.ResponseWriter, r *http.Request) {
	// TODO check if feed is already in OPML
	// adds a given URL to OPML
	var req AddFeedRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Sprintf("cannot parse json: %v\n", err.Error()))
	}
	url := req.url
	// from a URL, I need to get the title
	rss, err := ParseRSSFromURL(url)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("cannot parse RSS: %v\n", err.Error()))
	}
	title := rss.GetTitle()
	s.opml.AddFeed(title, url, "rss")
	w.WriteHeader(http.StatusOK)
}
