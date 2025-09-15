package domain

import (
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Feeds struct {
	ID      int
	Name    string
	Url     string
	Created time.Time
	Updated time.Time
}

type Articles struct {
	Title       string
	Link        string
	Description string
	Published   string
	Created     time.Time
	Updated     time.Time
	FeedID      int
}
