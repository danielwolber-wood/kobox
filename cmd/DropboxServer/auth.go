package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Code for non-PKCE Auth v1; no longer function because the TokenManager implementation was changed, but preserved for reference

// AuthFlow returns a RefreshToken using the non-PKCE authorization flow
func AuthFlow() (*RequestRefreshTokenOptions, error) {
	err := godotenv.Load(envFile)
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

func GetAuthCode(appKey string) string {
	//https://www.dropbox.com/oauth2/authorize?client_id=<APP_KEY>&response_type=code&redirect_uri=<REDIRECT_URI>&state=<STATE>
	authUrl := fmt.Sprintf("https://www.dropbox.com/oauth2/authorize?client_id=%s&token_access_type=offline&response_type=code", appKey)
	fmt.Printf("Navigate to \n%s\n and enter the key here: ", authUrl)
	var code string
	fmt.Scan(&code)
	return code
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

func (tm *TokenManager) refreshToken() error {

	opts := RequestRefreshTokenOptions{
		RefreshToken: tm.token.RefreshToken,
		ClientID:     tm.ClientID,
		//ClientSecret: tm.ClientSecret,
	}
	token, err := RequestRefreshToken(opts)
	if err != nil {
		return err
	}
	tm.token = *token
	tm.expiresAt = time.Now().Add(3 * time.Hour)
	return nil
}
