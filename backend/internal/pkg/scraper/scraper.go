package scraper

import (
	"github.com/cockroachdb/errors"
	"github.com/sejunpark/headline/backend/constant"
	"github.com/sejunpark/headline/backend/internal/pkg/model"
	"github.com/sejunpark/headline/backend/internal/pkg/scraper/builder"
	"github.com/sejunpark/headline/backend/internal/pkg/scraper/thebell"
	"log/slog"
	"time"
)

type Scraper struct {
	builder.Builder
}

func NewScraper(source constant.Source) (scraper *Scraper, cleanup func(), err error) {
	var b builder.Builder

	switch source {
	case constant.SourceThebell:
		b, cleanup, err = thebell.NewThebellScraperBuilder()
		return &Scraper{Builder: b}, cleanup, err
	default:
		return nil, nil, errors.AssertionFailedf("unsupported source: %v", source)
	}
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
