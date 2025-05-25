package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func (tm *TokenManager) GetValidToken() (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if tm.expiresAt.Before(time.Now()) {
		if err := tm.refreshTokenPKCE(); err != nil {
			return "", err
		}
	}
	return tm.token.AccessToken, nil
}

func (tm *TokenManager) refreshTokenPKCE() error {

	opts := RequestRefreshTokenPKCEOptions{
		RefreshToken: tm.token.RefreshToken,
		ClientID:     tm.ClientID,
	}
	token, err := RequestRefreshTokenPKCE(opts)
	if err != nil {
		return err
	}
	tm.token = *token
	tm.expiresAt = time.Now().Add(3 * time.Hour)
	return nil
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
