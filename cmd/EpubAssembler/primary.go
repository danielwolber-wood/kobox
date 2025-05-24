package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// Fetch accepts a URL string and returns HTML
func Fetch(url URL) (string, error) {
	resp, err := http.Get(string(url))
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
func Extract(worker JSWorker, html HTML) (ReadabilityObject, error) {
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
		return Epub{}, err
	}
	return epub, nil
}

// Upload accepts an UploadObject and returns nothing
func Upload(uploadObject UploadObject) error {
	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload", bytes.NewReader(uploadObject.Data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg", fmt.Sprintf("{\"path\": \"%s\"}", uploadObject.DestinationPath))

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("status: %v\n", resp.Status)
	fmt.Printf("body: %v\n", body)
	return nil
}
