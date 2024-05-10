package model

import (
	"fmt"
	"github.com/sejunpark/headline/backend/internal/pkg/util"
	"net/url"
	"time"
)

type Article struct {
	ArticleInfo
	Content string
}

func (a *Article) IsContentScraped() bool {
	return a.Content != ""
}

type ArticleInfo struct {
	Keywords        map[string]bool // set implementation of keywords used to search this article
	Title           string
	Summary         string
	CreatedDateTime time.Time
	UpdateDateTime  time.Time
	Url             *url.URL
	Source          string
	SourceUrl       *url.URL
}

func (a *ArticleInfo) IsValid() bool {
	keywordsNotEmpty := len(a.Keywords) > 0
	titleNotEmpty := a.Title != ""
	createdDateTimeNotEmpty := a.CreatedDateTime != time.Time{}
	urlIsValid := util.IsUrlValid(a.Url)
	sourceNotEmpty := a.Source != ""
	sourceUrlIsValid := util.IsUrlValid(a.SourceUrl)

	return keywordsNotEmpty && titleNotEmpty && createdDateTimeNotEmpty && urlIsValid && sourceNotEmpty && sourceUrlIsValid
}

func (a *ArticleInfo) String() string {
	return fmt.Sprintf("Title: %s\nSummary: %s\b", a.Title, a.Summary)
}
