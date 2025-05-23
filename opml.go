package main

import (
	"encoding/xml"
	"io"
	"os"
)

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

type Head struct {
	Title string `xml:"title"`
}

type Body struct {
	Outline Outline `xml:"outline"`
}

type Outline struct {
	Text     string    `xml:"text,attr"`
	Type     string    `xml:"type,attr,omitempty"`
	XMLUrl   string    `xml:"xmlUrl,attr,omitempty"`
	Outlines []Outline `xml:"outline,omitempty"`
}

func ParseOPML(xmlData string) (*OPML, error) {
	var opml OPML
	err := xml.Unmarshal([]byte(xmlData), &opml)
	if err != nil {
		return nil, err
	}
	return &opml, nil
}

func ParseOPMLFromReader(reader io.Reader) (*OPML, error) {
	var opml OPML
	decoder := xml.NewDecoder(reader)
	err := decoder.Decode(opml)
	if err != nil {
		return nil, err
	}
	return &opml, nil
}

func ExtractFeeds(outline Outline) []Outline {
	var feeds []Outline
	if outline.XMLUrl != "" {
		feeds = append(feeds, outline)
	}

	for _, nested := range outline.Outlines {
		feeds = append(feeds, ExtractFeeds(nested)...)
	}

	return feeds

}

func ParseOPMLFile(filename string) (*OPML, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ParseOPMLFromReader(file)
}
