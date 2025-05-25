package main

import (
	"fmt"
	"github.com/go-shiori/go-readability"
	"sync"

	//"golang.org/x/oauth2"
	"log"
	"time"
)

type Server struct {
	tokenManager    *TokenManager
	jsWorkerFactory *JSWorkerFactory
	jobQueue        chan Job
}

func newServer(opts RequestRefreshTokenOptions) (*Server, error) {

	jobQueue := make(chan Job, 256)
	jsWorkerFactory, err := NewJSWorkerFactory()
	if err != nil {
		return nil, fmt.Errorf("Could not create worker factory: %v\n", err)
	}
	token, err := RequestRefreshToken(opts)
	if err != nil {
		return nil, fmt.Errorf("could not get refresh token: %v\n", err)
	}
	tokenManager := TokenManager{mu: sync.RWMutex{}, token: *token, expiresAt: time.Now().Add(time.Second * 14000), ClientID: opts.ClientID, ClientSecret: opts.ClientSecret}
	return &Server{
		tokenManager:    &tokenManager,
		jobQueue:        jobQueue,
		jsWorkerFactory: jsWorkerFactory,
	}, nil
}

func (s *Server) worker(n int) {
	/* jsWorker, err := s.jsWorkerFactory.NewJSWorker()
	if err != nil {
		log.Printf("Failed to start worker: %v\n", err)
		return
	}
	*/
	for job := range s.jobQueue {
		switch job.currentStep {
		case StepPrefetch:
			// TODO check if a manual title was passed in and, if so, use that
			fmt.Println("fetching")
			article, err := readability.FromURL(string(job.url), 30*time.Second)
			if err != nil {
				log.Printf("Error fetching article: %v\n", err)
			}
			ro := ReadabilityObject{
				Title:   article.Title,
				Content: article.Content,
				Excerpt: article.Excerpt,
			}
			job.readabilityObject = ro
			job.currentStep = StepExtracted
			s.jobQueue <- job
		case StepExtracted:
			// run epub generation, ro -> epub []bytes
			fmt.Println("generating")
			epub, err := Generate(job.readabilityObject)
			if err != nil {
				log.Printf("Error generating epub: %v\n", err)
				continue
			}
			job.epub = epub
			job.currentStep = StepGenerated
			s.jobQueue <- job
		case StepGenerated:
			// construct upload object then upload to dropbox
			u := UploadObject{
				Data:            job.epub,
				Mimetype:        "application/epub+zip",
				DestinationPath: fmt.Sprintf("/Apps/Rakuten Kobo/%v.epub", job.readabilityObject.Title),
			}
			fmt.Println("uploading")
			fmt.Printf("u is %v\n", u)
			accessToken, err := s.tokenManager.GetValidToken()
			if err != nil {
				log.Printf("error getting new access token")
			}
			s.tokenManager.mu.Unlock()
			err = Upload(u, accessToken)
			if err != nil {
				log.Printf("Error uploading: %v\n", err)
			}
			fmt.Println("done")
			// TODO add to queue
			continue
		case StepUploaded:
			// do nothing now, but should be a "send success to client that made request" step
			continue
		default:
			continue
		}
	}
}
