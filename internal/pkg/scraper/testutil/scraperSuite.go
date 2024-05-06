package testutil

import (
	"github.com/sejunpark/headline/internal/pkg/scraper"
	"github.com/stretchr/testify/suite"
)

type BaseScraperSuite struct {
	suite.Suite
	Scraper *scraper.Scraper
	cleanup func()
}

func (ts *BaseScraperSuite) SetupBuilderSuite(builder scraper.Builder, cleanup func()) {
	s, cleanup, err := scraper.NewScraper(builder, cleanup)
	ts.NoErrorf(err, "failed to initialize Scraper: %v", err)
	ts.Scraper = s
	ts.cleanup = cleanup
}

func (ts *BaseScraperSuite) TearDownSuite() {
	ts.cleanup()
}
