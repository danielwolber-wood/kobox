package main

import (
	"os"
	"path/filepath"
)

type Server struct {
	readabilityParser *ReadabilityParser
	opml              *OPML
	opmlFilename      string
}

func newServer() (*Server, error) {
	appdata, err := GetAppDataDir("kobox-mono")
	if err != nil {
		return nil, err
	}
	opmlFilename := filepath.Join(appdata, "opml.xml")
	var opml OPML
	if !FileExists(opmlFilename) {
		_, err = os.Create(opmlFilename)
		if err != nil {
			return nil, err
		}
		opml = NewOPML()
		opml.Save(opmlFilename)
	} else {
		opml, err = ParseOPMLFile(opmlFilename)
		if err != nil {
			return nil, err
		}
	}
	parser, err := NewReadabilityParser()
	if err != nil {
		return nil, err
	}
	return &Server{
		readabilityParser: parser,
		opml:              &opml,
		opmlFilename:      opmlFilename,
	}, nil
}
