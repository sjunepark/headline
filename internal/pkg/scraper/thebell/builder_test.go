package thebell

import (
	"github.com/sejunpark/headline/internal/pkg/scraper/testutil"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"testing"
	"time"
)

type theBellScraperBuilderSuite struct {
	testutil.BaseScraperBuilderSuite
	logLevel slog.Level
}

func (ts *theBellScraperBuilderSuite) SetupSuite() {
	// Have to set log level in the beginning, not in the BaseSuite since this has to be set before any other process.
	ts.logLevel = slog.SetLogLoggerLevel(slog.LevelDebug)

	builder, cleanup, err := NewThebellScraperBuilder()
	ts.NoErrorf(err, "failed to initialize TheBellScraperBuilder: %v", err)
	ts.SetupScraperBuilderSuite(builder, cleanup)
}

func (ts *theBellScraperBuilderSuite) TearDownSuite() {
	ts.SetupTearDownSuite()
	slog.SetLogLoggerLevel(ts.logLevel)
}

func TestThebellScraperBuilderSuite(t *testing.T) {
	suite.Run(t, new(theBellScraperBuilderSuite))
}

func (ts *theBellScraperBuilderSuite) Test_fetchArticlesPage() {
	ts.Run("when keyword is valid it should return a page with newsList and paging", func() {
		p, err := ts.Builder.FetchArticlesPage("삼성전자", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		newsList, err := p.Element(".newsList")
		ts.NoError(err, "failed to get newsList")
		ts.NotNil(newsList)

		paging, err := p.Element(".paging")
		ts.NoError(err, "failed to get paging")
		ts.NotNil(paging)
	})

	ts.Run("when keyword is empty it should return an error", func() {
		p, err := ts.Builder.FetchArticlesPage("", time.Time{})
		ts.Error(err, "expected an error")
		ts.Nil(p, "expected nil")
	})
}

func (ts *theBellScraperBuilderSuite) Test_FetchNextPage() {
	ts.Run("when nextPage exists should return true and a different next page", func() {
		initialPage, err := ts.Builder.FetchArticlesPage("삼성전자", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		nextPage, exists := ts.Builder.FetchNextPage(initialPage)
		ts.True(exists)
		ts.NotNil(nextPage)

		ts.NotEqual(initialPage.Text(), nextPage.Text())
	})

	ts.Run("when nextPage does not exist should return false and error", func() {
		p, err := ts.Builder.FetchArticlesPage("청호나이스", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		nextPage, exists := ts.Builder.FetchNextPage(p)
		ts.False(exists)
		ts.Nil(nextPage)
	})
}

func (ts *theBellScraperBuilderSuite) Test_parseArticlesPage() {
	p, err := ts.Builder.FetchArticlesPage("삼성전자", time.Time{})
	ts.NoError(err, "failed to fetch articles page")

	ts.Run("when page has articles it should return ArticleInfos", func() {
		articleInfos, err := ts.Builder.ParseArticlesPage(p)
		ts.NoError(err, "failed to parse articles page")
		ts.NotEmpty(articleInfos)

		for _, articleInfo := range articleInfos {
			ts.Truef(articleInfo.IsValid(), "expected valid articleInfo, got %v", articleInfo)
		}
	})
}
