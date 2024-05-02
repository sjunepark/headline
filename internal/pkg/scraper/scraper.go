package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/model"
	"log/slog"
	"net/url"
	"sync"
)

// sourceScraper is an interface that defines the methods that a sourceScraper should implement.
//
// fetchUrlsToScrape should return a channel of urls to scrape, and it should handle pagination.
// It should also be thread-safe.
//
// fetchArticle should return an article for a given url, and should be thread-safe.
type sourceScraper interface {
	fetchUrlsToScrape(keyword string) (<-chan url.URL, error)
	fetchArticle(url.URL) (model.Article, error)
	cleanup()
	//	todo: implement context to cancel the scraping in certain conditions
}

type Scraper struct {
	sourceScraper
}

func NewScraper(s sourceScraper) (scraper *Scraper, cleanup func()) {
	scraper = &Scraper{s}
	return scraper, s.cleanup
}

func (s *Scraper) Scrape(keyword string) (<-chan model.Article, error) {
	urlsToScrape, err := s.fetchUrlsToScrape(keyword)
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
			article, articleFetchErr := s.fetchArticle(u)
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

func (s *Scraper) fetchAllUrlsToScrape(keywords []string) (<-chan url.URL, error) {
	urls := make(chan url.URL)
	defer close(urls)

	wg := sync.WaitGroup{}
	for _, keyword := range keywords {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			urlsToScrape, err := s.fetchUrlsToScrape(keyword)
			if err != nil {
				slog.Error("failed to fetch urls to scrape", "error", err)
				return
			}
			for u := range urlsToScrape {
				urls <- u
			}
		}(keyword)
	}

	wg.Wait()
	return urls, nil
}
