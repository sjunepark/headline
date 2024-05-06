package model

import (
	"github.com/sejunpark/headline/internal/pkg/util"
	"net/url"
	"time"
)

type Article struct {
	ArticleInfos
	Content string
}

func (a *Article) IsContentScraped() bool {
	return a.Content != ""
}

type ArticleInfos struct {
	Keywords        map[string]bool // set implementation of keywords used to search this article
	Title           string
	Summary         string
	CreatedDateTime time.Time
	UpdateDateTime  time.Time
	Url             *url.URL
	Source          string
	SourceUrl       *url.URL
}

func (a *ArticleInfos) IsValid() bool {
	keywordsNotEmpty := len(a.Keywords) > 0
	titleNotEmpty := a.Title != ""
	createdDateTimeNotEmpty := a.CreatedDateTime != time.Time{}
	urlIsValid := util.IsUrlValid(a.Url)
	sourceNotEmpty := a.Source != ""
	sourceUrlIsValid := util.IsUrlValid(a.SourceUrl)

	return keywordsNotEmpty && titleNotEmpty && createdDateTimeNotEmpty && urlIsValid && sourceNotEmpty && sourceUrlIsValid
}
