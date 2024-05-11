package thebell

import (
	"github.com/sejunpark/headline/backend/internal/pkg/scraper/testutil"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ThebellScraperSuite struct {
	testutil.BaseScraperSuite
}

func (ts *ThebellScraperSuite) SetupSuite() {
	builder, cleanup, err := NewThebellScraperBuilder()
	ts.NoErrorf(err, "failed to initialize TheBellScraperBuilder: %v", err)
	ts.SetupBuilderSuite(builder, cleanup)
}

func TestThebellScraperSuite(t *testing.T) {
	suite.Run(t, new(ThebellScraperSuite))
}

func (ts *ThebellScraperSuite) TestScrape() {
	articleInfos, err := ts.Scraper.Scrape("보령바이오파마", time.Time{})
	ts.NoError(err, "failed to scrape")
	ts.NotEmpty(articleInfos)
}
