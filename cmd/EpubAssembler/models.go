package main

import (
	_ "embed"
	"github.com/dop251/goja"
)

const (
	StepPrefetch = iota
	StepFetched
	StepExtracted
	StepGenerated
	StepUploaded
)

type Step byte

type Epub []byte

type URL string

type HTML string

type URLRequestObject struct {
	Url   URL    `json:"url"`
	Title string `json:"title"`
}

type HTMLRequestObject struct {
	Html  HTML   `json:"html"`
	Title string `json:"title"`
}

type ReadabilityObject struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Excerpt string `json:"excerpt"`
}

type UploadObject struct {
	Data            []byte
	Mimetype        string
	DestinationPath string
}

//go:embed readability.js
var readabilityJS string

type JSWorkerFactory struct {
	readabilityProgram *goja.Program
}

type JSWorker struct {
	vm *goja.Runtime
}

type Job struct {
	currentStep       Step
	url               URL
	fullText          HTML
	epub              Epub
	readabilityObject ReadabilityObject
	uploadObject      UploadObject
}
