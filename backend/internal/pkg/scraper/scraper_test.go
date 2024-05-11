package scraper

import (
	"github.com/sejunpark/headline/backend/constant"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ThebellScraperSuite struct {
	BaseScraperSuite
}

func (ts *ThebellScraperSuite) SetupSuite() {
	ts.SetupScraperSuite(constant.SourceThebell)
}

func TestThebellScraperSuite(t *testing.T) {
	suite.Run(t, new(ThebellScraperSuite))
}

func (ts *ThebellScraperSuite) TestScrape() {
	articleInfos, err := ts.Scraper.Scrape("보령바이오파마", time.Time{})
	ts.NoError(err, "failed to scrape")
	ts.NotEmpty(articleInfos)
}
