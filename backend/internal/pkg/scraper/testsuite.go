package scraper

import (
	"github.com/sejunpark/headline/backend/constant"
	"github.com/stretchr/testify/suite"
	"log/slog"
)

type BaseScraperSuite struct {
	suite.Suite
	Scraper  *Scraper
	cleanup  func()
	logLevel slog.Level
}

func (ts *BaseScraperSuite) SetupScraperSuite(source constant.Source) {
	ts.logLevel = slog.SetLogLoggerLevel(slog.LevelDebug)

	var err error
	ts.Scraper, ts.cleanup, err = NewScraper(source)
	ts.NoErrorf(err, "failed to initialize Scraper: %v", err)
}

func (ts *BaseScraperSuite) TearDownSuite() {
	ts.cleanup()
	slog.SetLogLoggerLevel(ts.logLevel)
}
