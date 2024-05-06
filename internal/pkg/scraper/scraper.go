package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/model"
	"time"
)

type Builder interface {
	FetchArticlesPage(keyword string, startDate time.Time) (*ArticlesPage, error)
	FetchNextPage(currentPage *ArticlesPage) (nextPage *ArticlesPage, exists bool)
	ParseArticlesPage(p *ArticlesPage) ([]*model.ArticleInfos, error)
}

type Scraper struct {
	Builder
}

func NewScraper(builder Builder, cleanup func()) (*Scraper, func(), error) {
	return &Scraper{Builder: builder}, cleanup, nil
}

func (s *Scraper) Scrape(keyword string, startDate time.Time) ([]*model.ArticleInfos, error) {
	var infos []*model.ArticleInfos
	var currentPage *ArticlesPage
	var nextPageExists bool
	var err error

	for {
		currentPage, err = s.FetchArticlesPage(keyword, startDate)
		if err != nil {
			return nil, err
		}
		currentInfos, parseErr := s.ParseArticlesPage(currentPage)
		if parseErr != nil {
			return nil, parseErr
		}
		infos = append(infos, currentInfos...)

		currentPage, nextPageExists = s.FetchNextPage(currentPage) //nolint:ineffassign,staticcheck
		if !nextPageExists {
			break
		}
	}

	return infos, nil
}
