package testutil

import (
	"github.com/sejunpark/headline/internal/pkg/scraper"
	"github.com/stretchr/testify/suite"
	"log/slog"
)

type BaseScraperSuite struct {
	suite.Suite
	Scraper  *scraper.Scraper
	cleanup  func()
	logLevel slog.Level
}

func (ts *BaseScraperSuite) SetupBuilderSuite(builder scraper.Builder, cleanup func()) {
	ts.logLevel = slog.SetLogLoggerLevel(slog.LevelDebug)
	s, cleanup, err := scraper.NewScraper(builder, cleanup)
	ts.NoErrorf(err, "failed to initialize Scraper: %v", err)
	ts.Scraper = s
	ts.cleanup = cleanup
}

func (ts *BaseScraperSuite) TearDownSuite() {
	ts.cleanup()
	slog.SetLogLoggerLevel(ts.logLevel)
}
