package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func RandomString(n int) (string, error) {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-._~")
	b := make([]rune, n)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		b[i] = chars[n.Int64()]
	}
	return string(b), nil
}

func HashString(s string) [32]byte {
	hash := sha256.Sum256([]byte(s))
	return hash
}

func GetPKCE() (*PKCECode, error) {
	codeVerifier, err := RandomString(128)
	if err != nil {
		return nil, err
	}
	hash := HashString(codeVerifier)
	codeChallenge := strings.TrimRight(base64.URLEncoding.EncodeToString(hash[:]), "=")
	pkce := PKCECode{
		CodeVerifier:  codeVerifier,
		CodeChallenge: codeChallenge,
	}
	return &pkce, nil
}

func GetAuthCodePKCE(appKey string, pkce *PKCECode) (string, error) {
	//https://www.dropbox.com/oauth2/authorize?client_id=<APP_KEY>&response_type=code&code_challenge=<CHALLENGE>&code_challenge_method=<METHOD>
	pkce, err := GetPKCE()
	challengeMethod := "S256"
	if err != nil {
		return "", err
	}
	authUrl := fmt.Sprintf("https://www.dropbox.com/oauth2/authorize?client_id=%s&token_access_type=offline&response_type=code&code_challenge=%s&code_challenge_method=%s", appKey, pkce.CodeChallenge, challengeMethod)
	fmt.Printf("Navigate to \n%s\n and enter the key here: ", authUrl)
	var code string
	fmt.Scan(&code)
	return code, nil

}

func RequestAccessTokenPKCE(opts RequestAccessTokenPKCEOptions) (*Token, error) {
	/* curl https://api.dropbox.com/oauth2/token \
	   -d code=<AUTHORIZATION_CODE> \
	   -d grant_type=authorization_code \
	   -d redirect_uri=<REDIRECT_URI> \
	   -d code_verifier=<VERIFICATION_CODE> \
	   -d client_id=<APP_KEY>
	*/
	data := url.Values{}
	data.Set("code", opts.AuthCode)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", opts.ClientID)
	data.Set("code_verifier", opts.CodeVerifier)
	resp, err := http.Post(
		"https://api.dropbox.com/oauth2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
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

func RequestRefreshTokenPKCE(opts RequestRefreshTokenPKCEOptions) (*Token, error) {
	/*
		curl https://api.dropbox.com/oauth2/token \
		    -d grant_type=refresh_token \
		    -d refresh_token=<REFRESH_TOKEN> \
		    -d client_id=<APP_KEY>
	*/
	data := url.Values{}
	data.Set("refresh_token", opts.RefreshToken)
	data.Set("client_id", opts.ClientID)
	resp, err := http.Post(
		"https://api.dropbox.com/oauth2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
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

func AuthFlowPKCE() (*RequestRefreshTokenPKCEOptions, error) {
	err := godotenv.Load(envFile)
	if err != nil {
		return nil, err
	}
	refreshToken := os.Getenv("REFRESH_TOKEN")
	if refreshToken == "" {
		// If there is no refresh token, go through auth process
		pkce, err := GetPKCE()
		if err != nil {
			return nil, err
		}
		authCode, err := GetAuthCodePKCE(dropboxClientId, pkce)
		if err != nil {
			return nil, err
		}
		opts := RequestAccessTokenPKCEOptions{
			AuthCode:     authCode,
			ClientID:     dropboxClientId,
			CodeVerifier: pkce.CodeVerifier,
		}
		accessToken, err := RequestAccessTokenPKCE(opts)
		if err != nil {
			return nil, err
		}
		refreshToken = accessToken.RefreshToken
		updateEnvFile("REFRESH_TOKEN", refreshToken)
	}

	return &RequestRefreshTokenPKCEOptions{
		RefreshToken: refreshToken,
		ClientID:     dropboxClientId,
	}, nil
}
