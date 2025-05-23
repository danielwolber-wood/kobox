// https://www.rssboard.org/rss-specification
package main

import (
	"encoding/xml"
	"io"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Copyright   string `xml:"copyright"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	PubDate     string    `xml:"pubDate"`
	Enclosure   Enclosure `xml:"enclosure"`
	Link        string    `xml:"link"`
	GUID        string    `xml:"guid"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

func ParseRSS(rssData string) (*RSS, error) {
	var rss RSS
	err := xml.Unmarshal([]byte(rssData), &rss)
	if err != nil {
		return nil, err
	}
	return &rss, nil
}

func ParseRSSFromReader(reader io.Reader) (*RSS, error) {
	var rss RSS
	decoder := xml.NewDecoder(reader)
	err := decoder.Decode(&rss)
	if err != nil {
		return nil, err
	}
	return &rss, nil
}

func ExtractItems(rss RSS) []Item {
	var items []Item
	for _, item := range rss.Channel.Items {
		items = append(items, item)
	}
	return items
}
