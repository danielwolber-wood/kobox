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
		// TODO add Task for taking as input a full HTML page
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
			ro := ReadabilityObject{
				Title:   article.Title,
				Content: article.Content,
				Excerpt: article.Excerpt,
			}
			//log.Printf("content is: \n%v\n", article.Content)
			//log.Printf("textcontent is: \n%v\n", article.TextContent)
			job.readabilityObject = ro
			job.taskType = TaskGenerate
			s.jobQueue <- job
		case TaskGenerate:
			// run epub generation, ro -> epub []bytes
			fmt.Println("generating")
			epub, err := Generate(job.readabilityObject)
			if err != nil {
				log.Printf("Error generating epub: %v\n", err)
				continue
			}
			job.epub = epub
			job.taskType = TaskUpload
			s.jobQueue <- job
		case TaskUpload:
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
