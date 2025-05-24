package main

import "github.com/dop251/goja"

type Epub []byte

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

type ReadabilityParser struct {
	vm *goja.Runtime
}
