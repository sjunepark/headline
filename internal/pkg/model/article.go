package model

import (
	"net/url"
	"time"
)

type Article struct {
	Title           string
	CreatedDateTime time.Time
	UpdatedDateTime time.Time
	Source          string
	Url             *url.URL
	Summary         string
	Content         string
	Keywords        map[string]bool
	ScrapeStatus    ScrapeStatus
}

// ScrapeStatus is an enum for the status of scraping process
// 1. UrlScraped: the url has been scraped, which would be from scraping the list of articles
// 2. ListItemParsed: the list item element has been parsed,
// which would usually have parsed the title, url, and datetime, etc.
// 3. ArticleParsed: the article has been parsed, which would have parsed the text content
type ScrapeStatus int

const (
	UrlScraped ScrapeStatus = iota
	ListItemParsed
	ArticleParsed
)

func (a *Article) IsUrlValid() bool {
	return a.Url.Scheme != "" && a.Url.Host != ""
}
