package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func GetAuthCode(appKey string) string {
	// todo add code validation
	//https://www.dropbox.com/oauth2/authorize?client_id=<APP_KEY>&response_type=code&redirect_uri=<REDIRECT_URI>&state=<STATE>
	authUrl := fmt.Sprintf("https://www.dropbox.com/oauth2/authorize?client_id=%s&token_access_type=offline&response_type=code", appKey)
	fmt.Printf("Navigate to \n%s\n and enter the key here: ", authUrl)
	var code string
	fmt.Scan(&code)
	return code
}

func RevokeToken(accessToken string) (*http.Response, error) {
	url := "https://api.dropboxapi.com/2/auth/token/revoke"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Token revoked successfully. Response: %s\n", string(body))
	return resp, nil
}

func ExchangeCodeForToken(opts TokenExchangeOptions) (*Token, error) {
	data := url.Values{}
	data.Set("code", opts.AuthCode)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", opts.ClientID)
	data.Set("client_secret", opts.ClientSecret)
	/*
		curl https://api.dropbox.com/oauth2/token \
		    -d code=<AUTHORIZATION_CODE> \
		    -d grant_type=authorization_code \
		    -d client_id=<APP_KEY> \
		    -d client_secret=<APP_SECRET>
	*/
	resp, err := http.Post(
		"https://api.dropbox.com/oauth2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	var token Token
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func RequestRefreshToken(opts RequestRefreshTokenOptions) (*Token, error) {
	data := url.Values{}
	data.Set("refresh_token", opts.RefreshToken)
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", opts.ClientID)
	data.Set("client_secret", opts.ClientSecret)
	resp, err := http.Post(
		"https://api.dropbox.com/oauth2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	var token Token
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// AuthFlow returns a RefreshToken
func AuthFlow() (*RequestRefreshTokenOptions, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("No .env file found or couldn't load it, using environment variables: %v\n", err)
	}
	refreshToken := os.Getenv("REFRESH_TOKEN")
	dropboxClientId := os.Getenv("DROPBOX_APP_KEY")
	if dropboxClientId == "" {
		return nil, fmt.Errorf("DROPBOX_APP_KEY not set")
	}

	dropboxClientSecret := os.Getenv("DROPBOX_APP_SECRET")
	if dropboxClientSecret == "" {
		return nil, fmt.Errorf("DROPBOX_APP_SECRET not set")
	}
	if refreshToken == "" {
		// If there is no refresh token, go through auth process
		authCode := GetAuthCode(dropboxClientId)
		tokenExchangeOptions := TokenExchangeOptions{
			AuthCode:     authCode,
			ClientID:     dropboxClientId,
			ClientSecret: dropboxClientSecret,
		}
		token, err := ExchangeCodeForToken(tokenExchangeOptions)
		if err != nil {
			return nil, err
		}
		refreshToken = token.RefreshToken
		updateEnvFile("REFRESH_TOKEN", refreshToken)
	}

	return &RequestRefreshTokenOptions{
		RefreshToken: refreshToken,
		ClientID:     dropboxClientId,
		ClientSecret: dropboxClientSecret,
	}, nil
}

func updateEnvFile(key, value string) error {
	envFile := ".env"
	content, err := os.ReadFile(envFile)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	updated := false
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, key) {
			lines[i] = fmt.Sprintf("%s=%s", key, value)
			updated = true
			break
		}
	}
	if !updated {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	return os.WriteFile(envFile, []byte(strings.Join(lines, "\n")), 0644)
}

func (tm *TokenManager) GetValidToken() (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if tm.expiresAt.Before(time.Now()) {
		if err := tm.refreshToken(); err != nil {
			return "", err
		}
	}
	return tm.token.AccessToken, nil
}

func (tm *TokenManager) refreshToken() error {

	opts := RequestRefreshTokenOptions{
		RefreshToken: tm.token.RefreshToken,
		ClientID:     tm.ClientID,
		ClientSecret: tm.ClientSecret,
	}
	token, err := RequestRefreshToken(opts)
	if err != nil {
		return err
	}
	tm.token = *token
	tm.expiresAt = time.Now().Add(3 * time.Hour)
	return nil
}
