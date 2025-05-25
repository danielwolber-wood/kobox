package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/danielwolber-wood/kobox/internal/response"
	"html/template"
	"io"
	"log"
	"net/http"
)

//go:embed static/config.html
var configHTML string

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

func (s *Server) handlerConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// user strings.Replace to insert correct auth link
		tmpl, err := template.New("config").Parse(configHTML)
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
		}
		// todo PKCE flow integration required
		pkce, err := GetPKCE()
		if err != nil {
			http.Error(w, "PKCE Error", http.StatusInternalServerError)
		}
		authUrl := fmt.Sprintf("https://www.dropbox.com/oauth2/authorize?client_id=%s&token_access_type=offline&response_type=code&code_challenge=%s&code_challenge_method=%s", dropboxClientId, pkce.CodeChallenge, "S256")
		data := ConfigData{
			AuthURL:       authUrl,
			CodeVerifier:  pkce.CodeVerifier,
			CodeChallenge: pkce.CodeChallenge,
		}
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, data)
		return
	}
	var authCode string
	if r.Method == "POST" {
		authCode = r.FormValue("auth_code")
		//codeChallenge := r.FormValue("code_challenge")
		codeVerifier := r.FormValue("code_verifier")
		if authCode == "" {
			http.Error(w, "Authorization code is required", http.StatusBadRequest)
			return
		}
		fmt.Printf("auth code is %s\n", authCode)
		opts := RequestAccessTokenPKCEOptions{
			AuthCode:     authCode,
			ClientID:     dropboxClientId,
			CodeVerifier: codeVerifier,
		}
		accessToken, err := RequestAccessTokenPKCE(opts)
		log.Printf("access token PKCE is %v\n", accessToken)
		if err != nil {
			http.Error(w, "Unable to get access token", http.StatusInternalServerError)
		}
		s.configureTokenManager(accessToken)

		fmt.Fprintf(w, `
            <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 50px auto; text-align: center;">
                <h2>Success!</h2>
                <p>Authorization code received: <code>%s</code></p>
                <p>You can now close this window.</p>
            </div>
        `, authCode)
	}
}
