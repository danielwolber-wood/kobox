package main

import (
	"embed"
	"github.com/dop251/goja"
)

const (
	StepPrefetch = iota
	StepFetched
	StepExtracted
	StepGenerate
	StepUpload
)

type Step byte

type Epub []byte

type URL string

type HTML string

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
	currentStep Step
	data        any
}
