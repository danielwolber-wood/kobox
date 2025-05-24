package main

import (
	"fmt"
	"log"
	"os"
)

type Server struct {
	dropboxApiKey   string
	jsWorkerFactory *JSWorkerFactory
	jobQueue        chan Job
}

func newServer() (*Server, error) {
	apiKey := os.Getenv("DROPBOX_API_KEY")
	if apiKey == "" {
		log.Fatal("DROPBOX_API_KEY not set")
	}
	jobQueue := make(chan Job, 256)
	jsWorkerFactory, err := NewJSWorkerFactory()
	if err != nil {
		log.Fatal("Could not create worker factory")
	}
	return &Server{
		dropboxApiKey:   apiKey,
		jobQueue:        jobQueue,
		jsWorkerFactory: jsWorkerFactory,
	}, nil
}

func (s *Server) worker(n int) {
	jsWorker, err := s.jsWorkerFactory.NewJSWorker()
	if err != nil {
		log.Printf("Failed to start worker: %v\n", err)
		return
	}
	for job := range s.jobQueue {
		switch job.currentStep {
		case StepPrefetch:
			fmt.Println("fetching")
			fullPage, err := Fetch(job.url)
			if err != nil {
				log.Printf("Error processing url: %v\n", err)
				continue
			}
			job.fullText = fullPage
			/*
				fmt.Println("extracting")
				ro, err := Extract(*jsWorker, job.fullText)
				if err != nil {
					log.Printf("Error extracting: %v\n", err)
					continue
				}
				job.readabilityObject = ro
				fmt.Printf("content is %v\n", job.readabilityObject.Content)
			*/
			epub, err := ConvertStringWithPandoc(fullPage, "html", "epub")
			if err != nil {
				log.Printf("Error converting with pandoc: %v\n", err)
				continue
			}
			title := job.readabilityObject.Title
			fmt.Printf("title is %v\n", title)
			ro := ReadabilityObject{Content: string(epub), Title: title}
			job.readabilityObject = ro
			fmt.Println("generating")
			/*epub, err = Generate(job.readabilityObject)
			if err != nil {
				log.Printf("Error generating epub: %v\n", err)
				continue
			} */
			job.epub = epub
			u := UploadObject{
				Data:            job.epub,
				Mimetype:        "application/epub+zip",
				DestinationPath: fmt.Sprintf("/Apps/Rakuten Kobo/%v.epub", job.readabilityObject.Title),
			}
			fmt.Println("uploading")
			fmt.Printf("u is %v\n", u)
			err = Upload(u, s.dropboxApiKey)
			if err != nil {
				log.Printf("Error uploading: %v\n", err)
			}
			fmt.Println("done")
		case StepFetched:
			// run extraction, html -> ro
			continue
		case StepExtracted:
			// run epub generation, ro -> epub []bytes
			continue
		case StepGenerated:
			// upload to dropbox, construct upload object
			continue
		case StepUploaded:
			// do nothing now, but should be a "send success to client that made request" step
			continue
		default:
			switch job.currentStep {
			case StepPrefetch:
				fmt.Println("fetching")
				fullPage, err := Fetch(job.url)
				if err != nil {
					log.Printf("Error processing url: %v\n", err)
					continue
				}
				job.fullText = fullPage
				fmt.Println("extracting")
				ro, err := Extract(*jsWorker, job.fullText)
				if err != nil {
					log.Printf("Error extracting: %v\n", err)
					continue
				}
				job.readabilityObject = ro
				fmt.Printf("content is %v\n", job.readabilityObject.Content)
				epub, err := ConvertStringWithPandoc(fullPage, "html", "epub")
				if err != nil {
					log.Printf("Error converting with pandoc: %v\n", err)
					continue
				}
				title := job.readabilityObject.Title
				fmt.Printf("title is %v\n", title)
				ro = ReadabilityObject{Content: string(epub), Title: title}
				job.readabilityObject = ro
				fmt.Println("generating")
				epub, err = Generate(job.readabilityObject)
				if err != nil {
					log.Printf("Error generating epub: %v\n", err)
					continue
				}
				job.epub = epub
				u := UploadObject{
					Data:            job.epub,
					Mimetype:        "application/epub+zip",
					DestinationPath: fmt.Sprintf("/Apps/Rakuten Kobo/%v.epub", job.readabilityObject.Title),
				}
				fmt.Println("uploading")
				fmt.Printf("u is %v\n", u)
				err = Upload(u, s.dropboxApiKey)
				if err != nil {
					log.Printf("Error uploading: %v\n", err)
				}
				fmt.Println("done")
				continue

			}
		}
	}
}
