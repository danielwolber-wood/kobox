package main

import (
	"io"
	"net/http"
)

// Fetch accepts a URL string and returns HTML
func Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Extract accepts HTML and returns a ReadabilityObject
func Extract(worker JSWorker, html string) (ReadabilityObject, error) {
	obj, err := worker.ParseHTML(html)
	if err != nil {
		return ReadabilityObject{}, err
	}
	return *obj, nil
}

// Generate accepts a ReadabilityObject and returns an Epub
func Generate(rr ReadabilityObject) (Epub, error) {
	html := GenerateHTML(rr.Title, rr.Title)
	epub, err := ConvertStringWithPandoc(html, "html", "epub")
	if err != nil {
		return err
	}
	return epub, nil
}

// Upload accepts an UploadObject and returns nothing
func Upload(UploadObject any) {}
