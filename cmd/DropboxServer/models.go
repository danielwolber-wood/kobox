package main

import (
	"io"
	"sync"
	"time"
)

const (
	TaskFetch = iota
	TaskExtract
	TaskGenerate
	TaskUpload
	TaskInform
)

type TaskType byte

type Epub []byte

type URL string

type HTML string

type TokenExchangeOptions struct {
	AuthCode     string
	ClientID     string
	ClientSecret string
}

type RequestRefreshTokenOptions struct {
	RefreshToken string
	ClientID     string
	ClientSecret string
}

type RequestRefreshTokenPKCEOptions struct {
	RefreshToken string
	ClientID     string
}

type RequestAccessTokenPKCEOptions struct {
	AuthCode     string
	ClientID     string
	CodeVerifier string
}

// NOTE access token is the short-lived one
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type TokenManager struct {
	mu        sync.RWMutex
	token     Token
	expiresAt time.Time
	ClientID  string
	//ClientSecret string
}

type URLRequestObject struct {
	Url   URL    `json:"url"`
	Title string `json:"title"`
}

type HTMLRequestObject struct {
	Html  HTML   `json:"html"`
	Title string `json:"title"`
}

type GenerateOptions struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Excerpt string `json:"excerpt"`
}

type UploadOptions struct {
	Data            []byte
	Mimetype        string
	DestinationPath string
}

type Job struct {
	taskType        TaskType
	url             URL
	htmlReader      io.Reader
	fullText        HTML
	epub            Epub
	generateOptions GenerateOptions
	uploadOptions   UploadOptions
}

type PKCECode struct {
	CodeVerifier  string
	CodeChallenge string
}
