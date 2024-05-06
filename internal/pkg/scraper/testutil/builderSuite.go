package testutil

import (
	"github.com/sejunpark/headline/internal/pkg/scraper"
	"github.com/stretchr/testify/suite"
)

type BaseScraperBuilderSuite struct {
	suite.Suite
	Builder scraper.Builder
	cleanup func()
}

func (ts *BaseScraperBuilderSuite) SetupScraperBuilderSuite(builder scraper.Builder, cleanup func()) {
	ts.Builder = builder
	ts.cleanup = cleanup
}

func (ts *BaseScraperBuilderSuite) TearDownSuite() {
	ts.cleanup()
}
