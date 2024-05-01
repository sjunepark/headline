package model

import "net/url"

type Article struct {
	Title           string
	CreatedDateTime string // format: "2006-01-02 15:04:05"
	UpdatedDateTime string // format: "2006-01-02 15:04:05"
	Source          string
	Url             url.URL
	Summary         string
	Content         string
	Keywords        map[string]bool
}
