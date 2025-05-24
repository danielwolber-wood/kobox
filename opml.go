// https://opml.org/spec2.opml
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

func NewOPML() *OPML {
	return &OPML{
		XMLName: xml.Name{},
		Version: "2.0",
		Head:    Head{Title: "Feed Subscriptions"},
		Body:    Body{Outline: Outline{}},
	}
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
	err := decoder.Decode(&opml)
	if err != nil {
		return nil, err
	}
	return &opml, nil
}

func ParseOPMLFile(filename string) (*OPML, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ParseOPMLFromReader(file)
}

func (opml *OPML) Save(filename string) error {
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	data, err := xml.Marshal(opml)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

func (outline *Outline) ExtractFeeds() []Outline {
	var feeds []Outline
	if outline.XMLUrl != "" {
		feeds = append(feeds, *outline)
	}

	for _, nested := range outline.Outlines {
		feeds = append(feeds, nested.ExtractFeeds()...)
	}

	return feeds
}

func (opml *OPML) AddFeedFromOutline(outline Outline) {
	opml.Body.Outline.Outlines = append(opml.Body.Outline.Outlines, outline)
}

func (opml *OPML) AddFeed(title, url, feedType string) {
	outline := Outline{Text: title, Type: feedType, XMLUrl: url, Outlines: []Outline{}}
	opml.AddFeedFromOutline(outline)
}
