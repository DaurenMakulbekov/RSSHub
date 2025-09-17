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

type Add struct {
	Name string
	Url  string
}

type SetInterval struct {
	Duration string
}

type SetWorkers struct {
	Count string
}

type List struct {
	Num string
}

type Delete struct {
	Name string
}

type ArticlesCommand struct {
	FeedName string
	Num      string
}

type Fetch struct{}

type Commands struct {
	Name string
	Add
	SetInterval
	SetWorkers
	List
	Delete
	ArticlesCommand
	Fetch
}
