package scraper

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ThebellScraperSuite struct {
	suite.Suite
	scraper        *Scraper
	cleanupScraper func()
	// cleanupScraper should clean up the ThebellScraper, which will be tested
	cleanupThebellScraper func()
}

func (ts *ThebellScraperSuite) SetupTest() {
	thebell, cleanupThebellScraper, err := NewThebellScraper()
	ts.NoErrorf(err, "failed to initialize ThebellScraper: %v", err)
	ts.cleanupThebellScraper = cleanupThebellScraper

	scraper, cleanupScraper := NewScraper(thebell)
	ts.scraper = scraper
	ts.cleanupScraper = cleanupScraper
}

func (ts *ThebellScraperSuite) TearDownTest() {
	ts.cleanupScraper()
}

func TestThebellScraperSuite(t *testing.T) {
	suite.Run(t, new(ThebellScraperSuite))
}

func (ts *ThebellSuite) TestThebellScraper_cleanup() {
	version, err := ts.scraper.browser.rodBrowser.Version()
	ts.NoError(err)
	ts.NotEmpty(version.Product)

}
