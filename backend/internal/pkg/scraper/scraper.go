package scraper

import (
	"github.com/cockroachdb/errors"
	"github.com/sejunpark/headline/backend/internal/pkg/model"
	"log/slog"
	"time"
)

type Scraper struct {
	Builder
}

func NewScraper(builder Builder, cleanup func()) (*Scraper, func(), error) {
	return &Scraper{Builder: builder}, cleanup, nil
}

// Builder is an interface that defines the methods which a scraper builder must implement.
//
// FetchArticlesPage should return an ArticlesPage object for the given keyword and start date.
// The implementation of ArticlesPage is up to the builder, and effects how other methods handle it.
//
// FetchNextPage should return the next ArticlesPage if it exists.
// It should also check if the current page and the next page are different.
// If not, return false.
// Since an error is not returned from this function, it should log all errors.
//
// ParseArticlesPage should parse and return model.ArticleInfo objects.
type Builder interface {
	FetchArticlesPage(keyword string, startDate time.Time) (*ArticlesPage, error)
	FetchNextPage(currentPage *ArticlesPage) (nextPage *ArticlesPage, err error)
	ParseArticlesPage(p *ArticlesPage) ([]*model.ArticleInfo, error)
}

func (s *Scraper) Scrape(keyword string, startDate time.Time) ([]*model.ArticleInfo, error) {
	functionName := "Scraper.Scrape"

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
		if currentPage == nil {
			return nil, errors.AssertionFailedf("currentPage is nil. currentPage=%v,infos=%v, nextPageExists=%v", currentPage, infos, nextPageExists)
		}
		currentPage, err = s.FetchNextPage(currentPage)
		if err != nil {
			slog.Error("failed to fetch next page.", "function", functionName, "error", err)
			break
		}
		currentInfos, parseErr := s.ParseArticlesPage(currentPage)
		if parseErr != nil {
			slog.Error("failed to parse articles page.", "function", functionName, "error", parseErr)
			break
		}
		infos = append(infos, currentInfos...)
		slog.Debug("appended ArticleInfos.", "function", functionName, "appendedCount", len(currentInfos), "totalCount", len(infos))
	}

	return infos, nil
}
