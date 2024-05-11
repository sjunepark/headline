package app

import (
	"context"
	"github.com/sejunpark/headline/backend/internal/pkg/model"
	"sync"
)

type ScrapeService struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var scrape *ScrapeService
var onceScraper sync.Once

func Scrape() *ScrapeService {
	if scrape == nil {
		onceScraper.Do(func() {
			scrape = &ScrapeService{}
		})
	}
	return scrape
}

func (s *ScrapeService) Start(ctx context.Context) {
	s.ctx, s.cancel = context.WithCancel(ctx)
}

func (s *ScrapeService) Scrape(source string, keyword string) model.ArticlesPage {
	// TODO: Implement this
	return model.ArticlesPage{}
}
