package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/model"
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"time"
)

type scraperBuilder interface {
	fetchArticlesPage(keyword string, startDate time.Time) (*rodext.Element, error)
	fetchNextPage(currentPage *rodext.Element) (page *rodext.Element, exists bool, err error)
	parseArticlesPage(p *rodext.Element) ([]*model.ArticleInfos, error)
}

type Scraper struct {
	scraperBuilder
}

func NewScraper(builder scraperBuilder, cleanup func()) (*Scraper, func(), error) {
	return &Scraper{scraperBuilder: builder}, cleanup, nil
}

func (s *Scraper) scrape(keyword string, startDate time.Time) ([]*model.ArticleInfos, error) {
	var infos []*model.ArticleInfos
	var currentPage *rodext.Element
	var nextPageExists bool
	var err error

	for {
		currentPage, err = s.fetchArticlesPage(keyword, startDate)
		if err != nil {
			return nil, err
		}
		currentInfos, parseErr := s.parseArticlesPage(currentPage)
		if parseErr != nil {
			return nil, parseErr
		}
		infos = append(infos, currentInfos...)

		currentPage, nextPageExists, err = s.fetchNextPage(currentPage) //nolint:ineffassign,staticcheck
		if err != nil {
			break
		}
		if !nextPageExists {
			break
		}
	}

	return infos, nil
}
