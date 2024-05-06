package testutil

import (
	"github.com/sejunpark/headline/internal/pkg/scraper"
	"github.com/stretchr/testify/suite"
	"log/slog"
)

type BaseScraperBuilderSuite struct {
	suite.Suite
	Builder  scraper.Builder
	cleanup  func()
	logLevel slog.Level
}

func (ts *BaseScraperBuilderSuite) SetupScraperBuilderSuite(builder scraper.Builder, cleanup func()) {
	ts.logLevel = slog.SetLogLoggerLevel(slog.LevelDebug)

	ts.Builder = builder
	ts.cleanup = cleanup
}

func (ts *BaseScraperBuilderSuite) TearDownSuite() {
	slog.SetLogLoggerLevel(ts.logLevel)
	ts.cleanup()
}
