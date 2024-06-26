package thebell

import (
	"github.com/sejunpark/headline/backend/internal/pkg/scraper/testutil"
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
	ts.Run("when keyword is valid it should have an articles element and a paging element", func() {
		p, err := ts.Builder.FetchArticlesPage("삼성전자", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		ts.NotNil(p.Articles)
		ts.NotNil(p.PageNav)
	})

	ts.Run("when keyword is empty it should return an error", func() {
		p, err := ts.Builder.FetchArticlesPage("", time.Time{})
		ts.Error(err, "expected an error")
		ts.Nil(p, "expected nil")
	})
}

func (ts *theBellScraperBuilderSuite) Test_FetchNextPage() {
	ts.Run("when nextPage exists it should not return an error, and text should be different", func() {
		initialPage, err := ts.Builder.FetchArticlesPage("삼성전자", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		nextPage, err := ts.Builder.FetchNextPage(initialPage)
		ts.NoError(err, "failed to fetch next page")
		ts.NotNil(nextPage)

		ts.NotEqual(initialPage.Text(), nextPage.Text())
	})

	ts.Run("when nextPage does not exist should return false and error", func() {
		p, err := ts.Builder.FetchArticlesPage("청호나이스", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		nextPage, err := ts.Builder.FetchNextPage(p)
		ts.Error(err, "expected an error")
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
