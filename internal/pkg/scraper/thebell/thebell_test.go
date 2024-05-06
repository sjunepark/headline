package thebell

import (
	"github.com/sejunpark/headline/internal/pkg/scraper"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ThebellScraperSuite struct {
	suite.Suite
	scraper *scraper.Scraper
	cleanup func()
}

func (ts *ThebellScraperSuite) SetupSuite() {
	builder, cleanup, err := NewThebellScraperBuilder()
	ts.NoErrorf(err, "failed to initialize TheBellScraperBuilder: %v", err)
	s, cleanup, err := scraper.NewScraper(builder, cleanup)
	ts.NoErrorf(err, "failed to initialize Scraper: %v", err)
	ts.scraper = s
	ts.cleanup = cleanup
}

func (ts *ThebellScraperSuite) TearDownSuite() {
	ts.cleanup()
}

func (ts *ThebellScraperSuite) TestScrape() {
	articleInfos, err := ts.scraper.Scrape("보령바이오파마", time.Time{})
	ts.NoError(err, "failed to scrape")
	ts.NotEmpty(articleInfos)
	for _, articleInfo := range articleInfos {
		ts.Truef(articleInfo.IsValid(), "expected valid articleInfo, got %v", articleInfo)
		ts.T().Logf("articleInfo: %v", articleInfo)
	}

}

func TestThebellScraperSuite(t *testing.T) {
	suite.Run(t, new(ThebellScraperSuite))
}

type ThebellScraperBuilderSuite struct {
	suite.Suite
	builder *ScraperBuilder
	cleanup func()
}

func (ts *ThebellScraperBuilderSuite) SetupSuite() {
	builder, cleanup, err := NewThebellScraperBuilder()
	ts.NoErrorf(err, "failed to initialize TheBellScraperBuilder: %v", err)
	ts.builder = builder
	ts.cleanup = cleanup
}

func (ts *ThebellScraperBuilderSuite) TearDownSuite() {
	ts.cleanup()
}

func TestThebellBuilderSuite(t *testing.T) {
	suite.Run(t, new(ThebellScraperBuilderSuite))
}

func (ts *ThebellScraperBuilderSuite) Test_fetchArticlesPage() {
	ts.Run("when keyword is valid it should return a page with newsList and paging", func() {
		p, err := ts.builder.FetchArticlesPage("삼성전자", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		newsList, err := p.Element(".newsList")
		ts.NoError(err, "failed to get newsList")
		ts.NotNil(newsList)

		paging, err := p.Element(".paging")
		ts.NoError(err, "failed to get paging")
		ts.NotNil(paging)
	})

	ts.Run("when keyword is empty it should return an error", func() {
		p, err := ts.builder.FetchArticlesPage("", time.Time{})
		ts.Error(err, "expected an error")
		ts.Nil(p, "expected nil")
	})
}

func (ts *ThebellScraperBuilderSuite) Test_FetchNextPage() {
	ts.Run("when nextPage exits should return true and a different next page", func() {
		initialPage, err := ts.builder.FetchArticlesPage("삼성전자", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		nextPage, exists := ts.builder.FetchNextPage(initialPage)
		ts.True(exists)
		ts.NotNil(nextPage)

		ts.NotEqual(initialPage.Text(), nextPage.Text())
	})

	ts.Run("when nextPage does not exist should return false and error", func() {
		p, err := ts.builder.FetchArticlesPage("청호나이스", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		nextPage, exists := ts.builder.FetchNextPage(p)
		ts.False(exists)
		ts.Nil(nextPage)
	})
}

func (ts *ThebellScraperBuilderSuite) Test_parseArticlesPage() {
	p, err := ts.builder.FetchArticlesPage("삼성전자", time.Time{})
	ts.NoError(err, "failed to fetch articles page")

	ts.Run("when page has articles it should return ArticleInfos", func() {
		articleInfos, err := ts.builder.ParseArticlesPage(p)
		ts.NoError(err, "failed to parse articles page")
		ts.NotEmpty(articleInfos)

		for _, articleInfo := range articleInfos {
			ts.Truef(articleInfo.IsValid(), "expected valid articleInfo, got %v", articleInfo)
		}
	})
}
