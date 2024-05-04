package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/model"
	"github.com/sejunpark/headline/internal/pkg/rodext"
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
	fetchArticleList(keyword string) (<-chan *model.Article, error)
	fetchArticle(*url.URL) (*model.Article, error)
	parseArticleListItem(el *rodext.Element) (*model.Article, error)
	parseArticle(el *rodext.Element) (*model.Article, error)
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

func (s *Scraper) ScrapeAll(keyword []string) (<-chan *model.Article, error) {
	// todo: implement
	return nil, nil
}

func (s *Scraper) scrape(keyword string) (<-chan *model.Article, error) {
	articlesToScrape, err := s.fetchArticleList(keyword)
	if err != nil {
		return nil, err
	}

	articles := make(chan *model.Article)
	defer close(articles)

	wg := sync.WaitGroup{}
	for articleToScrape := range articlesToScrape {
		if !articleToScrape.IsUrlValid() {
			slog.Error("invalid url", "url", articleToScrape.Url.String())
			continue
		}
		u := articleToScrape.Url

		wg.Add(1)
		go func(u *url.URL) {
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

func (s *Scraper) fetchAllArticlesToScrape(keywords []string) (<-chan *model.Article, error) {
	articles := make(chan *model.Article)
	defer close(articles)

	wg := sync.WaitGroup{}
	for _, keyword := range keywords {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()

			articlesToScrape, err := s.fetchArticleList(keyword)
			if err != nil {
				slog.Error("failed to fetch articles to scrape", "error", err)
				return
			}
			for articleToScrape := range articlesToScrape {
				articles <- articleToScrape
			}
		}(keyword)
	}

	wg.Wait()
	return articles, nil
}
