package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/model"
	"time"
)

// sourceScraper
//
// fetchArticles should return a channel of ArticleMetadata for a given keyword.
// It should close the returning channel when there are no more articles to scrape.
// When no specific startDate is specified, a zero time.Time will be passed, which should be interpreted as "all time".
//
// fetchArticle should return an article for a given url, and should be thread-safe.
//
// cleanup should clean up any resources used by the sourceScraper, such as closing the browser.
type sourceScraper interface {
	cleanup()
	fetchArticles(keyword string, startDate time.Time) (<-chan *model.ArticleMetadata, error)
	String() string
}

type Scraper struct {
	sourceScraper
}

func NewScraper(s sourceScraper) (scraper *Scraper, cleanup func()) {
	scraper = &Scraper{s}
	return scraper, s.cleanup
}

// scrape fetches articles for the given keywords and start date.
// It returns a channel of articles, which will be closed when there are no more articles to scrape.
func (s *Scraper) scrape(keyword string, startDate time.Time) (<-chan *model.ArticleMetadata, error) {
	articles, err := s.fetchArticles(keyword, startDate)
	if err != nil {
		return nil, err
	}
	return articles, nil
}
