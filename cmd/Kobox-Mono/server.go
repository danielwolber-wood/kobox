package main

import (
	"fmt"
	"github.com/go-shiori/go-readability"
	"log"
	"s
	"strings"
	"log"
	"time"
)

type Server struct {
	tokenManager    *TokenManager
	jobQueue        chan Job
	tokenConfigured bool
}

func newServer() (*Server, error) {
	jobQueue := make(chan Job, 256)
	return &Server{
		//tokenManager: &tokenManager,
		jobQueue:        jobQueue,
		tokenConfigured: false,
	}, nil
}

func (s *Server) configureTokenManager(accessToken *Token) {
	token := Token{
		AccessToken:  accessToken.AccessToken,
		TokenType:    accessToken.TokenType,
		RefreshToken: accessToken.RefreshToken,
	}
	tokenManager := TokenManager{mu: sync.RWMutex{}, token: token, expiresAt: time.Now().Add(time.Second * 14000), ClientID: dropboxClientId}
	s.tokenManager = &tokenManager
	s.tokenConfigured = true
	s.tokenManager.GetValidToken()
}

func (s *Server) worker(n int) {
	for job := range s.jobQueue {
		switch job.taskType {
		case TaskFetch:
			// TODO implement a queue system of some kind for tasks that failed
			// TODO check if a manual title was passed in and, if so, use that
			fmt.Println("fetching")
			article, err := readability.FromURL(string(job.url), 30*time.Second)
			if err != nil {
				log.Printf("Error fetching article: %v\n", err)
				continue
			}
			ro := GenerateOptions{
				Title:   article.Title,
				Content: article.Content,
				Excerpt: article.Excerpt,
			}
			job.generateOptions = ro
			job.taskType = TaskGenerate
			s.jobQueue <- job
		case TaskExtract:
			// html request -> generateOptions
			title := job.generateOptions.Title
			fmt.Printf("title is %v\n", title)
			generateOption, err := Extract(job.htmlReader)
			if err != nil {
				log.Printf("error extracting article from reader: %v\n", err)
				continue
			}
			job.generateOptions = generateOption
			job.taskType = TaskGenerate
			s.jobQueue <- job
		case TaskGenerate:
			// run epub generation, generateOptions -> epub []bytes
			fmt.Println("generating")
			epub, err := Generate(job.generateOptions)
			if err != nil {
				log.Printf("Error generating epub: %v\n", err)
				continue
			}
			job.epub = epub
			job.taskType = TaskUpload
			s.jobQueue <- job
		case TaskUpload:
			// construct upload object then upload to dropbox
			u := UploadOptions{
				Data:            job.epub,
				Mimetype:        "application/epub+zip",
				DestinationPath: fmt.Sprintf("/Apps/Rakuten Kobo/%v.epub", sanitizeString(job.generateOptions.Title)),
			}
			fmt.Println("uploading")
			accessToken, err := s.tokenManager.GetValidToken()
			if err != nil {
				log.Printf("error getting new access token")
				continue
			}
			err = Upload(u, accessToken)
			if err != nil {
				log.Printf("Error uploading: %v\n", err)
				continue
			}
			fmt.Println("done")
			// TODO add to queue
			continue
		case TaskInform:
			// do nothing now, but should be a "send success to client that made request" step
			continue
		default:
			continue
		}
	}
}
func sanitizeString(input string) string {
	replacer := strings.NewReplacer(
		"?", "_",
		"\"", "_",
		"*", "_",
		"\\", "_",
		"|", "_",
		"/", "_",
	)
	return replacer.Replace(input)
}
