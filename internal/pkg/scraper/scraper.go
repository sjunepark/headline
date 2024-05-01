package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/model"
	"log/slog"
	"net/url"
	"sync"
)

// sourceScraper is an interface that defines the methods that a scraper should implement.
//
// fetchUrlsToScrape should return a channel of urls to scrape, and it should handle pagination.
// It should also be thread-safe.
//
// fetchArticle should return an article for a given url, and should be thread-safe.
type sourceScraper interface {
	fetchUrlsToScrape() (<-chan url.URL, error)
	fetchArticle(url.URL) (model.Article, error)
	//	todo: implement context to cancel the scraping in certain conditions
}

type Scraper struct {
	sourceScraper sourceScraper
}

func (s *Scraper) Scrape(keywords []string) (<-chan model.Article, error) {
	urlsToScrape, err := s.sourceScraper.fetchUrlsToScrape()
	if err != nil {
		return nil, err
	}

	articles := make(chan model.Article)
	defer close(articles)

	wg := sync.WaitGroup{}
	for u := range urlsToScrape {
		wg.Add(1)
		go func(u url.URL) {
			defer wg.Done()
			article, articleFetchErr := s.sourceScraper.fetchArticle(u)
			if articleFetchErr != nil {
				slog.Error("failed to fetch article", "url", u.String(), "error", articleFetchErr)
				return
			}
			articles <- article
		}(u)
	}

	wg.Wait()
	return articles, nil
}
