package main

import (
	"bytes"
	"fmt"
	"github.com/go-shiori/go-readability"
	"io"
	"net/http"
	"net/url"
)

// Extract accepts io.Reader and returns a GenerateOptions
func Extract(r io.Reader) (GenerateOptions, error) {
	urlObj, err := url.Parse("example.com")
	if err != nil {
		return GenerateOptions{}, err
	}
	article, err := readability.FromReader(r, urlObj)
	if err != nil {
		return GenerateOptions{}, err
	}
	return GenerateOptions{Title: article.Title, Content: article.Content, Excerpt: article.Excerpt}, nil
}

// Generate accepts a GenerateOptions and returns an Epub
func Generate(rr GenerateOptions) (Epub, error) {
	html := GenerateHTML(rr.Title, rr.Content)
	epub, err := ConvertStringWithPandoc(html, "html", "epub")
	if err != nil {
		return Epub{}, err
	}
	return epub, nil
}

// Upload accepts an UploadOptions and returns nothing
func Upload(uploadObject UploadOptions, accessToken string) error {
	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload", bytes.NewReader(uploadObject.Data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg", fmt.Sprintf("{\"path\": \"%s\"}", uploadObject.DestinationPath))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	fmt.Println(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("status: %v\n", resp.Status)
	fmt.Printf("body: %v\n", body)
	return nil
}
