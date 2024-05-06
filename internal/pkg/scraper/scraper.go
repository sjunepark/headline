package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/model"
	"log/slog"
	"time"
)

// Builder is an interface that defines the methods which a scraper builder must implement.
//
// FetchArticlesPage should return an ArticlesPage object for the given keyword and start date.
// The implementation of ArticlesPage is up to the builder, and effects how other methods handle it.
//
// FetchNextPage should return the next ArticlesPage if it exists.
// It should also check if the current page and the next page are different.
// Since an error is not returned, it should log all errors.
//
// ParseArticlesPage should parse and return model.ArticleInfo objects.
type Builder interface {
	FetchArticlesPage(keyword string, startDate time.Time) (*ArticlesPage, error)
	FetchNextPage(currentPage *ArticlesPage) (nextPage *ArticlesPage, exists bool)
	ParseArticlesPage(p *ArticlesPage) ([]*model.ArticleInfo, error)
}

type Scraper struct {
	Builder
}

func NewScraper(builder Builder, cleanup func()) (*Scraper, func(), error) {
	return &Scraper{Builder: builder}, cleanup, nil
}

func (s *Scraper) Scrape(keyword string, startDate time.Time) ([]*model.ArticleInfo, error) {
	currentPage, err := s.FetchArticlesPage(keyword, startDate)
	if err != nil {
		return nil, err
	}
	infos, err := s.ParseArticlesPage(currentPage)
	if err != nil {
		return nil, err
	}

	var nextPageExists bool
	for {
		currentPage, nextPageExists = s.FetchNextPage(currentPage)
		if !nextPageExists {
			slog.Debug("no more pages")
			break
		}
		currentInfos, parseErr := s.ParseArticlesPage(currentPage)
		if parseErr != nil {
			slog.Error("failed to parse articles page", "error", parseErr)
			break
		}
		infos = append(infos, currentInfos...)
	}

	return infos, nil
}
